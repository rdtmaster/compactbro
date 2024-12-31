package main

import (
	"github.com/rdtmaster/go-reddit/v4/reddit"
	"fmt"
	"sort"
	"time"
	"net/http"
	"github.com/labstack/echo/v4"
	"strings"
)

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

	return getOverview(c,
		c.Param("username"),
		c.QueryParam("after"),
		c.Param("page"),
		c.QueryParam("sort"),
		"overview_pt")
}
