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
	"runtime"
	"sort"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/kirsle/configdir"

	"strings"
	"time"

	"github.com/eknkc/amber"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rdtmaster/go-reddit/v4/reddit"
)

const (
	curVersion  = "CompactBro v0.80"
	kindComment = "t1"
	kindPost    = "t3"
)

type commentSubmP struct {
	Thing_id string `json:"thing_id"`
	Text     string `json:"text"`
}

type overviewWrap struct {
	PageTitle string
	Username  string
	Page      string
	After     string
	Sorting   string
	Items     []PostOrComment
	BodyClass string
}
type SubredditResponseWrapper struct {
	PageTitle string
	Sub       string
	After     string
	Sorting   string
	Items     []*reddit.Post
	BodyClass string
}
type PostOrComment struct {
	Kind string
	P    *reddit.Post
	C    *reddit.Comment
}
type voteResult struct {
	Direction string `json:"direction"`
	Thing_id  string `json:"thing_id"`
}
type CompactConfig struct {
	EcoMode            bool
	NightMode          bool
	DisplayFlairEmojis bool
	DefaultLimit       int
	LocalAddress       string
	HTTPS              struct {
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

func cVersion() (s string) {

	s = fmt.Sprintf("%s (%s/%s) @ %s (PID %d)", curVersion, runtime.GOOS, runtime.GOARCH, config.LocalAddress, os.Getpid())
	if config.HTTPS.Use {
		s += " HTTPS+"
	}

	return
}
func cssTheme() string {
	if config.NightMode {
		return "/static/night.css"
	}
	return "/static/compact.css"
}
func strToInt(s string) (res int) {
	if len(s) == 0 {
		return 0
	}
	res, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return
}

type DataWraper[T any] struct {
	PageTitle string
	Items     []*T
}
type PCWrapper struct {
	PageTitle string
	Items     []*reddit.Comment

	BodyClass string
	WP        *reddit.Post
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

const emptyHTML = template.HTML("")

func emoji(f reddit.RichFlair) template.HTML {
	if config.DisplayFlairEmojis {
		return template.HTML(fmt.Sprintf(`<span class="flairemoji" title="%s" style="background-image: url('%s')"></span>`, html.UnescapeString(f.A), html.UnescapeString(f.U)))
	}
	return emptyHTML

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

	for _, c := range replies.Comments {
		if c.HasMore() {
			c.Body_html += "<h1>I have more replies"
		}
	}

	var buf bytes.Buffer
	err := tpls["childComments"].Execute(&buf, replies)
	if err != nil {
		fmt.Println("=====Error in replies! ", err.Error())
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
	amber.FuncMap["cVersion"] = cVersion
	amber.FuncMap["emoji"] = emoji
	amber.FuncMap["cssTheme"] = cssTheme
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
	server.GET("/", frontpage)
	server.GET("/pt/", frontpagePT)
	server.GET("/r/:sub", subDefault)
	server.GET("/r/:sub/", subDefault)
	server.GET("/r/:sub/:sorting", subSorted)
	server.GET("/r/:sub/:sorting/", subSorted)

	server.GET("/pt/r/:sub/", subPT)
	server.GET("/r/:sub/comments/:id/:permalink/", submission)
	server.GET("/user/:sub/comments/:id/:permalink/", submission)
	server.GET("/u/:username/", overview)
	server.GET("/u/:username", overview)
	server.GET("/u/:username/:page/", overview)
	server.GET("/pt/u/:username/:page/", overviewPT)
	server.POST("/edit/", editThing)
	server.POST("/comment*", submitComment)
	server.GET("/vote/:direction/:thing_id/", vote)
	server.GET("/r/:sub/comments/:postID/:permalink/:commentID/", commentThread)

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

// comment thread display
// /r/:sub/comments/:postID/:permalink/:commentID/
func commentThread(c echo.Context) error {
	//TODO: <- move it to the library
	path := fmt.Sprintf("comments/%s/%s/%s", c.Param("postID"), c.Param("permalink"), c.Param("commentID"))
	req, err := client.NewRequest(http.MethodGet, path, nil)
	fmt.Println("=======", path, " ", req)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	pc := new(reddit.PostAndComments)
	_, err = client.Do(c.Request().Context(), req, pc)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	}
	pw := PCWrapper{
		PageTitle: pc.Post.Title,
		Items:     pc.Comments,
		WP:        pc.Post,

		BodyClass: "thread",
	}

	err = tpls["post"].Execute(c.Response(), pw)

	return err
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

func getSubreddit(c echo.Context, sub, after, sorting string, fp bool, tpl string) error {
	l := &reddit.ListOptions{
		Limit: config.DefaultLimit,
		After: after,
	}

	var pageTitle string
	var f func(context.Context, string, *reddit.ListOptions) ([]*reddit.Post, *reddit.Response, error)
	var sorting_ string
	if len(sorting) > 0 {
		sorting_ = sorting
	} else {
		sorting_ = "new"
	}
	switch strings.ToLower(sorting_) {
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
	posts, resp, err := f(c.Request().Context(), c.Param("sub"), l)

	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !config.EcoMode && !fp {
		sr, _, err := client.Subreddit.Get(c.Request().Context(), c.Param("sub"))
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		pageTitle = sr.Title
	} else if fp {
		pageTitle = "reddit - the front page of the internet"
	} else {
		if len(posts) > 0 {
			pageTitle = posts[0].SubredditNamePrefixed
		} else {
			pageTitle = "/r/" + sub
		}
	}

	start := time.Now()

	pw := SubredditResponseWrapper{
		PageTitle: pageTitle,
		After:     resp.After,
		Sub:       sub,
		Sorting:   sorting_,
		Items:     posts,
		BodyClass: "sub",
	}
	err = tpls[tpl].Execute(c.Response(), pw)

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
func frontpage(c echo.Context) error {
	return getSubreddit(c,
		"",
		c.QueryParam("after"),
		c.QueryParam("sort"),
		true,
		"frontpage")
}

func frontpagePT(c echo.Context) error {
	return getSubreddit(c,
		"",
		c.QueryParam("after"),
		c.QueryParam("sort"),
		true,
		"sub_pt")
}

// sub
func subDefault(c echo.Context) error {
	return getSubreddit(c,
		c.Param("sub"),
		c.QueryParam("after"),
		"new",
		false,
		"sub")
}

// /r/<sub>/pt/?sorting={hot|new...}&after=t****
func subPT(c echo.Context) error {
	return getSubreddit(c,
		c.Param("sub"),
		c.QueryParam("after"),
		c.QueryParam("sort"),
		false,
		"sub_pt")
}
func subSorted(c echo.Context) error {
	return getSubreddit(c,
		c.Param("sub"),
		c.QueryParam("after"),
		c.Param("sorting"),
		false,
		"sub")

}

// View post
func submission(c echo.Context) error {

	pc, _, err := client.Post.Get(c.Request().Context(), c.Param("id"))
	start := time.Now()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	pw := PCWrapper{
		PageTitle: pc.Post.Title,
		Items:     pc.Comments,
		WP:        pc.Post,

		BodyClass: "post",
	}

	err = tpls["post"].Execute(c.Response(), pw)

	elapsed := time.Since(start)
	fmt.Println("--------")
	fmt.Printf("comments rendered in %s", elapsed)
	fmt.Println()
	fmt.Println("--------")

	return err
}

// This function must be reworked or retested, hard to believe it actually works
func getOverview(c echo.Context, username, after, page, sorting, tpl string) error {

	l := &reddit.ListUserOverviewOptions{
		Sort: sorting,
	}
	l.After = after
	l.Limit = config.DefaultLimit
	var posts []*reddit.Post
	var comments []*reddit.Comment
	var resp *reddit.Response
	var err error
	switch page {
	case "submitted":
		posts, resp, err = client.User.PostsOf(c.Request().Context(), username, l)
		comments = nil
	case "comments":
		comments, resp, err = client.User.CommentsOf(c.Request().Context(), username, l)
		posts = nil
	default:
		posts, comments, resp, err = client.User.OverviewOf(c.Request().Context(), username, l)
	}
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println("==== after", after)
	a := make([]PostOrComment, len(posts)+len(comments))
	i := 0
	for _, post := range posts {
		a[i] = PostOrComment{
			Kind: "post",
			P:    post,
			C:    nil,
		}
		i++
	}

	for _, comment := range comments {
		a[i] = PostOrComment{
			Kind: "comment",
			P:    nil,
			C:    comment,
		}
		i++
	}
	if page == "overview" {
		sort.Slice(a, func(i, j int) bool {
			if a[i].Kind == "post" &&
				(a[i].P.Pinned || a[i].P.Stickied) {
				return true
			} else if a[j].Kind == "post" &&
				(a[j].P.Pinned || a[j].P.Stickied) {
				return false
			}

			switch sorting {
			case "hot", "top", "controversial": // TODO: <- create sorting by upvote ratio etc
				var t, u int
				if a[i].Kind == "post" {
					t = a[i].P.Score
				} else {
					t = a[i].C.Score
				}
				if a[j].Kind == "post" {
					u = a[j].P.Score
				} else {
					u = a[j].C.Score
				}
				return u < t

			default:
				var t, u time.Time
				if a[i].Kind == "post" {
					t = a[i].P.Created.Time
				} else {
					t = a[i].C.Created.Time
				}
				if a[j].Kind == "post" {
					u = a[j].P.Created.Time
				} else {
					u = a[j].C.Created.Time
				}
				return u.Before(t)
			}
		})
	}
	err = tpls[tpl].Execute(c.Response(), overviewWrap{
		PageTitle: "Overview for " + c.Param("username"),
		Username:  c.Param("username"),
		Sorting:   sorting,
		Page:      page,
		After:     resp.After,
		Items:     a,
		BodyClass: "overview",
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return err
}

// user overview
func overview(c echo.Context) error {
	sorting := strings.ToLower(c.QueryParam("sort"))
	if len(sorting) == 0 {
		sorting = "new"
	}
	page := strings.ToLower(c.Param("page"))
	if len(page) == 0 {
		page = "overview"
	}
	return getOverview(c,
		c.Param("username"),
		c.QueryParam("after"),
		page,
		sorting,
		"overview")
}

// /u/<user>/{hot|new|top|controversial}/pt/
func overviewPT(c echo.Context) error {
	type partRespOverview struct {
	}
	return getOverview(c,
		c.Param("username"),
		c.QueryParam("after"),
		c.Param("page"),
		c.QueryParam("sort"),
		"overview_pt")
}
