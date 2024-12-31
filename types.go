package main
import (
	"github.com/rdtmaster/go-reddit/v4/reddit"
	"github.com/eknkc/amber"
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
}
type SubredditResponseWrapper struct {
	PageTitle string
	Sub       string
	After     string
	Sorting   string
	Items     []*reddit.Post
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
	EcoMode              bool
	NightMode            bool
	DisplayFlairEmojis   bool
	DefaultLimit         int
	LocalAddress         string
	MarkMsgsUnreadOnView bool
	CheckMsgs            bool
	HTTPS                struct {
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

type DataWraper[T any] struct {
	PageTitle string
	Items     []*T
}
type PCWrapper struct {
	PageTitle string
	Items     []*reddit.Comment
	WP        *reddit.Post
	Thread    bool
}

type editP struct {
	Thing_id string `json:"thing_id"`
	Selftext string `json:"selftext"`
}
type editResult struct {
	Body string `json:"body"`
}

type userAttrs struct {
	Submitter bool
	Moderator bool
	Admin     bool
	Letters   string
}

type commentWrapper struct {
	Comment *reddit.Comment
}

type MessagesWrapper struct {
		After     string
		Page      string
		PageTitle string
		Sorting   string
		Items     []*reddit.Message
}

type partRespOverview struct {}

