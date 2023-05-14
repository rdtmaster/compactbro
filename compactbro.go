package main

import (
	"bytes"
	"context"
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

const (
	kindComment = "t1"
	kindPost    = "t3"
)

type commentSubmP struct {
	Thing_id string `json:"thing_id"`
	Text     string `json:"text"`
}

type voteResult struct {
	Direction string `json:"direction"`
	Thing_id  string `json:"thing_id"`
}
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

var oneItem = &reddit.ListOptions{Limit: 1}
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

type editP struct {
	Thing_id string `json:"thing_id"`
	Selftext string `json:"selftext"`
}
type editResult struct {
	Body string `json:"body"`
}

var tpls map[string]*template.Template

var client *reddit.Client

func emoji(f reddit.RichFlair) template.HTML {
	return template.HTML(fmt.Sprintf(`<span class="flairemoji" title="%s" style="background-image: url('%s')"></span>`, html.UnescapeString(f.A), html.UnescapeString(f.U)))
}
func getThumb(preview reddit.RedditPreview) string {
	if len(preview.Images) == 0 {
		return ""
	}
	return html.UnescapeString(preview.Images[0].Resolutions[0].URL)

}

func strNotEmpty(s string) bool {
	return len(s) > 0
}

func hasReplies(comment *reddit.Comment) bool {
	return len(comment.Replies.Comments) > 0
}

type userAttrs struct {
	Submitter bool
	Moderator bool
	Admin     bool
	Letters   string
}

func numTrues(bs ...bool) (n int) {
	n = 0
	for _, b := range bs {
		if b {
			n++
		}
	}
	return
}
func getDistinguished(distinguished string, isSubmitter bool) (ua userAttrs) {

	ua = userAttrs{
		Submitter: isSubmitter,
		Moderator: strings.Contains(distinguished, "moderator"),
		Admin:     strings.Contains(distinguished, "admin"),
	}
	ln := numTrues(ua.Moderator, ua.Admin, ua.Submitter)
	if ln == 0 {
		ua.Letters = ""
		return
	}
	tmp := make([]string, ln)
	i := 0
	if ua.Moderator {
		tmp[i] = "<a href=\"#\" class=\"moderator\" title=\"moderator of this subreddit, speaking officially\">M</a>"
		i++
	}
	if ua.Admin {
		tmp[i] = "<a href=\"#\" class=\"admin\" title=\"Reddit Administrator\">A</a>"
		i++
	}
	if ua.Submitter {
		tmp[i] = "<a href=\"#\" class=\"submitter\" title=\"submitter\">S</a>"
		i++
	}
	ua.Letters = "[" + strings.Join(tmp, ",") + "]"
	return
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
	amber.FuncMap["emoji"] = emoji

	amber.FuncMap["getThumb"] = getThumb
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
	server.File("/favicon.ico", "static/favicon.ico")

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
	server.GET("/r/:sub", subDefault)
	server.GET("/r/:sub/", subDefault)
	server.GET("/r/:sub/:sorting", subDisplay)
	server.GET("/r/:sub/:sorting/", subDisplay)
	server.GET("/r/:sub/comments/:id/:permalink/", submission)
	server.POST("/edit/", editThing)
	server.POST("/comment*", submitComment)
	server.GET("/vote/:direction/:thing_id/", vote)

	server.HEAD("/checkunread/", checkUnread)
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

func checkUnread(c echo.Context) error {
	if config.EcoMode {
		return c.NoContent(http.StatusNoContent)
	}
	ms, cs, _, err := client.Message.InboxUnread(c.Request().Context(), oneItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	if len(ms) > 0 || len(cs) > 0 {
		return c.NoContent(http.StatusOK)
	}
	return c.NoContent(http.StatusNoContent)
}

// up/down/remove-vote for post/coment
// /vote/{up|down|remove}/<thing_id>/
func vote(c echo.Context) error {
	fmt.Println("-----------")
	fmt.Println("Voting")
	fmt.Println("-----------")
	direction := c.Param("direction")
	thing_id := c.Param("thing_id")

	var f func(ctx context.Context, id string) (*reddit.Response, error)
	f = client.Post.Upvote
	switch direction {
	case "up":
		f = client.Post.Upvote
	case "down":
		f = client.Post.Downvote
	default:
		f = client.Post.RemoveVote
	}
	resp, err := f(c.Request().Context(), thing_id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(resp.StatusCode,
		voteResult{Direction: direction, Thing_id: thing_id})

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
func editThing(c echo.Context) error {
	var p editP
	err := c.Bind(&p)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	var r editResult
	if strings.HasPrefix(p.Thing_id, kindPost) {
		post, _, err := client.Post.Edit(c.Request().Context(), p.Thing_id, p.Selftext)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		r = editResult{post.Selftext_html}
	} else if strings.HasPrefix(p.Thing_id, kindComment) {
		comment, _, err := client.Comment.Edit(c.Request().Context(), p.Thing_id, p.Selftext)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		r = editResult{comment.Body_html}
	}

	return c.JSON(http.StatusOK, r)
}

func subSorted(c echo.Context, sorting string) error {
	var pageTitle string
	var f func(context.Context, string, *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error)
	switch strings.ToLower(sorting) {
	case "new":
		f = client.Subreddit.NewPosts

	case "hot":
		f = client.Subreddit.HotPosts
	case "rising":
		f = client.Subreddit.RisingPosts
	case "controversial":
		f = client.Subreddit.ControversialPosts
	case "top":
		f = client.Subreddit.TopPosts
	default:
		f = client.Subreddit.HotPosts
	}
	posts, _, err := f(c.Request().Context(), c.Param("sub"), nil)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !config.EcoMode {
		sr, _, err := client.Subreddit.Get(c.Request().Context(), c.Param("sub"))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pageTitle = sr.Title
	} else {
		if len(posts) > 0 {
			pageTitle = posts[0].SubredditNamePrefixed
		} else {
			pageTitle = "/r/" + c.Param("sub")
		}
	}

	start := time.Now()

	pw := DataWraper[reddit.Post]{
		PageTitle: pageTitle,
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

// sub
func subDefault(c echo.Context) error {
	return subSorted(c, "hot")
}
func subDisplay(c echo.Context) error {
	return subSorted(c, c.Param("sorting"))
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
