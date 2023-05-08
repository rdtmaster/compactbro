package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/kirsle/configdir"

	"strings"
	"time"

	"github.com/eknkc/amber"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rdtmaster/go-reddit/v3/reddit"
)

type CompactConfig struct {
	EcoMode      bool
	LocalAddress string
	HTTPS        struct {
		Use          bool
		LocalAddress string
		KeyPath      string
		CRTPath      string
	}
	Logging bool
	Auth    struct {
		Use      bool
		Username string
		Password string
	}
	Credentials     reddit.Credentials
	TemplateOptions amber.Options
}

var config CompactConfig
var server *echo.Echo

type DataWraper[T any] struct {
	PageTitle string
	Items     []*T
}
type PCWrapper struct {
	DataWraper[reddit.Comment]
	WP *reddit.Post
}

type postEditP struct {
	Thing_id string `json:"thing_id"`
	Selftext string `json:"selftext"`
}

var tpls map[string]*template.Template

var client *reddit.Client

func strNotEmpty(s string) bool {
	return len(s) > 0
}

func hasReplies(comment *reddit.Comment) bool {
	return len(comment.Replies.Comments) > 0
}

type userAttr struct {
	Class   string
	Letters string
}

var emptyAttrs = userAttr{"", ""}

func getDistinguished(distinguished string, isSubmitter bool) userAttr {
	if !isSubmitter && len(distinguished) == 0 {
		return emptyAttrs
	}
	class := distinguished
	if isSubmitter {
		if len(class) == 0 {
			class = "submitter"
		} else {
			class += " submitter"
		}
	}
	ln := strings.Count(class, " ") + 1

	tmp := make([]string, ln)
	i := 0
	if strings.Contains(class, "moderator") {
		tmp[i] = "<a href=\"#\" class=\"moderator\" title=\"moderator of this subreddit, speaking officially\">M</a>"
		i++
	}
	if strings.Contains(class, "admin") {
		tmp[i] = "<a href=\"#\" class=\"admin\" title=\"Reddit Administrator\">A</a>"
		i++
	}
	if strings.Contains(class, "submitter") {
		tmp[i] = "<a href=\"#\" class=\"submitter\" title=\"submitter\">S</a>"
		i++
	}
	return userAttr{
		Class:   class,
		Letters: "[" + strings.Join(tmp, ",") + "]",
	}

}
func processReplies(replies reddit.Replies) string {

	var buf bytes.Buffer
	err := tpls["childComments"].Execute(&buf, replies)
	if err != nil {
		fmt.Println("Error! ", err.Error())
		return err.Error()
	}

	return buf.String()
}
func isMine(author string) bool {
	return strings.
		EqualFold(client.Username, author)
}

func likesInt(l *bool) (likes int) {
	likes = 0
	if l != nil {
		if *l == true {
			likes = 1

		} else {
			likes = -1
		}
	}
	return
}

// Credit: https://www.socketloop.com/tutorials/golang-get-time-duration-in-year-month-week-or-day
func roundTime(input float64) int {
	var result float64
	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	i, _ := math.Modf(result)

	return int(i)
}

func calcrDate(diff time.Duration) (dVal int, dName string) {

	dName = ""
	dVal = 0
	secs := diff.Seconds()
	years := roundTime(secs / 31207680)
	if years > 0 {
		dVal = years
		dName = "year"
		return
	}

	months := roundTime(secs / 2600640)
	if months > 0 {
		dVal = months
		dName = "month"
		return
	}
	weeks := roundTime(secs / 604800)
	if weeks > 0 {
		dVal = weeks
		dName = "week"
		return
	}
	days := roundTime(secs / 86400)
	if days > 0 {
		dVal = days
		dName = "day"
		return
	}
	hours := roundTime(diff.Hours())
	if hours > 0 {
		dVal = hours
		dName = "hour"
		return
	}
	minutes := roundTime(diff.Minutes())
	if minutes > 0 {
		dVal = minutes
		dName = "minute"
		return
	}
	return
}
func dateAgo(t *reddit.Timestamp) string {
	dVal, dName := calcrDate(time.Since(t.Time))
	switch dVal {
	case 0:
		return "just now"
	case 1:
		return fmt.Sprintf("%d %s ago", dVal, dName)
	default:
		return fmt.Sprintf("%d %ss ago", dVal, dName)
	}
}

func main() {

	configPath := configdir.LocalConfig("compactbro")
	err := configdir.MakePath(configPath) // Ensure it exists.
	if err != nil {
		fmt.Println(err)
		return

	}
	configFile := filepath.Join(configPath, "compactbro.toml")
	fmt.Println("using config file ", configFile)

	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	amber.FuncMap["getDistinguished"] = getDistinguished
	amber.FuncMap["isMine"] = isMine

	amber.FuncMap["html"] = func(s string) template.HTML {
		return template.HTML(html.UnescapeString(s))
	}
	amber.FuncMap["dateAgo"] = dateAgo
	amber.FuncMap["likesInt"] = likesInt
	amber.FuncMap["strNotEmpty"] = strNotEmpty
	amber.FuncMap["hasReplies"] = hasReplies
	amber.FuncMap["processReplies"] = processReplies
	tpls, err = amber.CompileDir("templates",
		amber.DirOptions{Ext: ".amber", Recursive: true},
		config.TemplateOptions)

	if err != nil {
		fmt.Println("Error compiling templates ", err.Error())
		return
	}

	client, err = reddit.NewClient(config.Credentials)
	if err != nil {
		fmt.Println("Error Initializing client ", err.Error())
		return
	}
	server = echo.New()
	server.Static("/static", "static")

	// Middleware
	if config.Logging {
		server.Use(middleware.Logger())
	}
	server.Use(middleware.Recover())

	if config.Auth.Use {
		server.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			return strings.EqualFold(username, config.Auth.Username) && password == config.Auth.Password, nil
		}))
	}

	// Routes
	server.GET("/stop*", shutdown)
	server.GET("/r/:sub/", sub)
	server.GET("/r/:sub/comments/:id/:permalink/", submission)
	server.POST("/post/edit*", editPost)
	server.POST("/comment*", submitComment)

	// Start server
	go func() {
		if err := server.Start(config.LocalAddress); err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()
	if config.HTTPS.Use {
		if err := server.StartTLS(config.HTTPS.LocalAddress,
			config.HTTPS.CRTPath,
			config.HTTPS.KeyPath); err != http.ErrServerClosed {
			fmt.Println(err)
			os.Exit(0)

		}
	}
	//server.Logger.Fatal(server.Start(config.LocalAddress))

}

// Shut the server down
func shutdown(c echo.Context) error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		server.Close()
		os.Exit(0)
	}()
	return c.HTML(200, `<html><head><title>Goodbye</title><body bgcolor="#000123" text="#cdedfe"><h1 style="font-family: tahoma;right: 50%;bottom: 50%;transform: translate(50%,50%);position: absolute">Stopping server...</h1></body></html>`)
}
func indent(n int, s string) (res string) {
	spaces := "\n"
	for i := 1; i <= n; i++ {
		spaces += "    "
	}
	res = strings.ReplaceAll(s, "\n", spaces)
	return
}

type commentSubmP struct {
	Thing_id string `json:"thing_id"`
	Text     string `json:"text"`
}

// Submit comment
func submitComment(c echo.Context) error {
	var co commentSubmP
	err := c.Bind(&co)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	comment, _, err := client.Comment.Submit(c.Request().Context(), co.Thing_id, co.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	type commentWrapper struct {
		Comment *reddit.Comment
	}
	return tpls["oneComment"].Execute(c.Response(), commentWrapper{comment})

}

// Edit post
func editPost(c echo.Context) error {
	var p postEditP
	err := c.Bind(&p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	post, resp, err := client.Post.Edit(c.Request().Context(), p.Thing_id, p.Selftext)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(resp.StatusCode, post)

}

// sub
func sub(c echo.Context) error {

	sr, _, err := client.Subreddit.Get(c.Request().Context(), c.Param("sub"))
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	posts, _, err := client.Subreddit.HotPosts(c.Request().Context(), c.Param("sub"), nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	start := time.Now()

	pw := DataWraper[reddit.Post]{
		PageTitle: sr.Title,
		Items:     posts,
	}

	err = tpls["sub"].Execute(c.Response(), pw)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	elapsed := time.Since(start)
	fmt.Println("--------")
	fmt.Printf("sub rendered in %s", elapsed)
	fmt.Println()
	fmt.Println("--------")
	return err

}

// View post
func submission(c echo.Context) error {

	pc, _, err := client.Post.Get(c.Request().Context(), c.Param("id"))
	start := time.Now()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	pw := struct {
		PageTitle string
		Items     []*reddit.Comment
		WP        *reddit.Post
	}{}
	pw.PageTitle = pc.Post.Title
	pw.Items = pc.Comments
	pw.WP = pc.Post

	err = tpls["post"].Execute(c.Response(), pw)

	elapsed := time.Since(start)
	fmt.Println("--------")
	fmt.Printf("comments rendered in %s", elapsed)
	fmt.Println()
	fmt.Println("--------")

	return err
}
