package main

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"context"
	"github.com/labstack/echo/v4"
	"github.com/rdtmaster/go-reddit/v4/reddit"
	
)
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
