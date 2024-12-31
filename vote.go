package main

import(
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rdtmaster/go-reddit/v4/reddit"
	"net/http"
	"context"
)

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
