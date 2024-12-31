package main

import(
	"fmt"
	"net/http"
	"strings"
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/rdtmaster/go-reddit/v4/reddit"
)

func cleanCommentID(id string) (cleanID string) {
	cleanID, _ = strings.CutPrefix(id, kindComment+"_")
	return
}
func hasReplies(comment *reddit.Comment) bool {
	return len(comment.Replies.Comments) > 0
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



// comment thread display
// /r/:sub/comments/:postID/:permalink/:commentID/
func commentThread(c echo.Context) error {
	//TODO: <- move it to the library
	path := fmt.Sprintf("comments/%s/%s/%s", c.Param("postID"), c.Param("permalink"), c.Param("commentID"))
	req, err := client.NewRequest(http.MethodGet, path, nil)

	fmt.Println("====ua ", req.UserAgent())
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
		Thread:    true,
	}

	err = tpls["post"].Execute(c.Response(), pw)

	return err
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
	
	return tpls["oneComment"].Execute(c.Response(), commentWrapper{comment})

}
