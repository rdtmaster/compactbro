package main

import(
	"github.com/labstack/echo/v4"
	
	"net/http"
	"time"
	"fmt"
	"strings"
	
)

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
		Thread:    false,
	}

	err = tpls["post"].Execute(c.Response(), pw)

	elapsed := time.Since(start)
	fmt.Println("--------")
	fmt.Printf("comments rendered in %s", elapsed)
	fmt.Println()
	fmt.Println("--------")

	return err
}
