package main

import(
	"fmt"
	"context"
	"sort"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/rdtmaster/go-reddit/v4/reddit"
	
)

func sentWrapped(ctx context.Context, opts *reddit.ListOptions) ([]*reddit.Message, []*reddit.Message, *reddit.Response, error) {
	messages, resp, err := client.Message.Sent(ctx, opts)
	return messages, nil, resp, err
}


func checkUnread(c echo.Context) error {
	if !config.CheckMsgs {
		return c.NoContent(http.StatusNoContent)
	}
	oneItem := &reddit.ListOptions{
		Limit: 1,
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

func getMessages(c echo.Context, page, after, tpl string) error {
	l := &reddit.ListOptions{
		After: after,
		Limit: config.DefaultLimit,
	}

	var page_ string
	if len(page) > 0 {
		page_ = page
	} else {
		page_ = "all"
	}
	var f func(context.Context, *reddit.ListOptions) ([]*reddit.Message, []*reddit.Message, *reddit.Response, error)
	switch page_ {
	case "all":
		f = client.Message.Inbox
	case "unread":
		f = client.Message.InboxUnread
	case "comments":
		f = client.Message.InboxComments
	case "selfreply":
		f = client.Message.InboxSelfReplies
	case "mentions":
		f = client.Message.InboxMentions
	case "messages":
		f = client.Message.InboxMessages
	case "sent":
		f = sentWrapped

	default:
		f = client.Message.Inbox
	}
	comments, messages, resp, err := f(c.Request().Context(), l)
	fmt.Println("=====messages======")
	unread := make([]string, 0)
	a := make([]*reddit.Message, len(messages)+len(comments))
	index := 0
	for _, msg := range messages {
		a[index] = msg
		index++
		if msg.New {
			unread = append(unread, msg.FullID)
		}
	}
	for _, msg := range comments {
		a[index] = msg
		index++
		if msg.New {
			unread = append(unread, msg.FullID)
		}
	}

	if len(unread) > 0 && config.MarkMsgsUnreadOnView {
		go func() {
			_, err = client.Message.Read(context.Background(), unread...)
			if err != nil {
				fmt.Println("Marking unread err: ", err.Error())
			} else {
				fmt.Println("Successfully marked ", len(unread), " messages as read.")
			}
		}()
	}
	if len(messages) > 0 && len(comments) > 0 {
		sort.Slice(a, func(i, j int) bool {
			return a[j].Created.Time.Before(a[i].Created.Time)
		})
	}

	
	wm := &MessagesWrapper{
		After:     resp.After,
		Page:      page_,
		PageTitle: "messages: " + page,
		Sorting:   "",
		Items:     a,
	}
	err = tpls[tpl].Execute(c.Response(), wm)
	return err
}
func messageInbox(c echo.Context) error {
	return getMessages(c,
		c.Param("page"),
		c.QueryParam("after"),
		"messages")

}
func messageInboxPT(c echo.Context) error {
	return getMessages(c,
		c.Param("page"),
		c.QueryParam("after"),
		"messages_pt")

}
