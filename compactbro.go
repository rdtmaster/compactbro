package main

import (
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
	EcoMode         bool
	Credentials     reddit.Credentials
	TemplateOptions amber.Options
}

var config CompactConfig
var server *echo.Echo

type DataWraper[T any] struct {
	Title string
	Items []*T
}
type PCWrapper struct {
	DataWraper[reddit.Comment]
	WP *PostWrapper
}
type PostWrapper struct {
	IsDistinguished bool
	IsMine          bool
	HasLinkFlair    bool
	DateAgo         string
	Post            *reddit.Post
}

type postEditP struct {
	Thing_id string `json:"thing_id"`
	Selftext string `json:"selftext"`
}

var tpls map[string]*template.Template

var client *reddit.Client

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
func rDate(t *reddit.Timestamp) string {
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
func wrapPost(post *reddit.Post) *PostWrapper {

	return &PostWrapper{
		IsMine:          strings.EqualFold(strings.ToLower(client.Username), strings.ToLower(post.Author)),
		IsDistinguished: len(post.Distinguished) > 0,
		HasLinkFlair:    len(post.LinkFlairText) > 0,
		Post:            post,
		DateAgo:         rDate(post.Created),
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

	amber.FuncMap["html"] = func(s string) template.HTML {
		return template.HTML(html.UnescapeString(s))
	}
	amber.FuncMap["likesInt"] = likesInt
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
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// Routes
	server.GET("/stop", shutdown)
	server.GET("/r/:sub", sub)
	server.GET("/r/:sub/comments/:id/:permalink", submission)
	server.POST("/post/edit", editPost)
	server.POST("/comment", submitComment)

	// Start server
	server.Logger.Fatal(server.Start(":80"))
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

func processComments(cs []*reddit.Comment) (s string) {

	for _, c := range cs {
		entryclass := "unvoted"
		upclass := "arrow up login-required "
		downclass := "arrow down login-required"
		if c.Likes != nil {
			if *(c.Likes) == true {
				entryclass = "likes"
				upclass += " upmod"

			} else {
				entryclass = "dislikes"
				downclass += " downmod"
			}
		}

		midcol := fmt.Sprintf(`	<div class="midcol"><div class="%s"></div><div class="%s"></div></div>
`, upclass, downclass)
		authorclass := "author may-blank id-" + c.FullID
		userattrs := ""
		if c.IsSubmitter {
			userattrs = "[<a class=\"submitter\" title=\"submitter\" href=\"#\">S</a>]"
			authorclass += " submitter"
		}
		entry := fmt.Sprintf(`	<div class="entry %s">
		<div class="tagline">
			<a class="%s" href="/u/%[3]s">%[3]s</a>
			<span class="userattrs">%s</span> <span class="score dislikes">%d points</span><span class="score unvoted">%d points</span><span class="score likes">%d points</span>  %s
		</div>
	</div>
	<a href="" class="options_link"></a>
	<form action="#" class="usertext">
		<input name="thing_id" type="hidden" value="%s"/>
		<div class="usertext-body">%s</div>
	</form>
	<div class="clear options_expando hidden"></div>
`, entryclass, authorclass, c.Author, userattrs,
			c.Score-1, c.Score, c.Score+1,
			rDate(c.Created), c.FullID, html.UnescapeString(c.Body_html))
		s += fmt.Sprintf(`
<div class ="thing comment" data-id="%s">
%s %s
<div class="commentspacer"></div>

`, c.FullID, midcol, entry)

		if len(c.Replies.Comments) > 0 {
			s += "<div class=\"child\">" + processComments(c.Replies.Comments) + "</div>"
		}
		s += "</div>"
	}
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
	comment, resp, err := client.Comment.Submit(c.Request().Context(), co.Thing_id, co.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.HTML(resp.StatusCode, processComments([]*reddit.Comment{comment}))
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

	mps := make([]*PostWrapper, len(posts), len(posts))
	for i, post := range posts {
		mps[i] = wrapPost(post)
	}

	pw := DataWraper[PostWrapper]{
		Title: sr.Title,
		Items: mps,
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
	start := time.Now()
	pc, _, err := client.Post.Get(c.Request().Context(), c.Param("id"))

	t := template.HTML(processComments(pc.Comments))
	fmt.Println("-------------")
	fmt.Println(t)
	fmt.Println("-------------")

	pw := struct {
		Title string
		Items []*reddit.Comment
		WP    *PostWrapper
		CBody template.HTML
	}{}
	pw.Title = pc.Post.Title
	pw.Items = pc.Comments
	pw.WP = wrapPost(pc.Post)
	pw.CBody = t

	err = tpls["post"].Execute(c.Response(), pw)

	elapsed := time.Since(start)
	fmt.Println("--------")
	fmt.Printf("comments rendered in %s", elapsed)
	fmt.Println()
	fmt.Println("--------")

	return err
}
