package main

import(
	"fmt"
	"strings"
	"runtime"
	"html"
	"html/template"
	"os"
	"strconv"
	"net/url"
	"github.com/rdtmaster/go-reddit/v4/reddit"
)
func linkFromContext(context string) (link string) {
	tmp, _ := strings.CutPrefix(context, "/r/")
	as := strings.Split(tmp, "/")
	link = fmt.Sprintf("/r/%s/comments/%s/", as[0], as[2])
	return

}
func cleanLink(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return link
	}
	s := u.Path
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	if !strings.HasSuffix(s, "/") {
		s = s + "/"
	}
	return s
}
func isPostID(parent string) bool {
	return strings.HasPrefix(parent, kindPost)
}
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


func emoji(f reddit.RichFlair) template.HTML {
	if config.DisplayFlairEmojis {
		return template.HTML(fmt.Sprintf(`<span class="flairemoji" title="%s" style="background-image: url('%s')"></span>`, html.UnescapeString(f.A), html.UnescapeString(f.U)))
	}
	return emptyHTML

}
func getThumb(preview reddit.RedditPreview) string {
	if len(preview.Images) <= 0 || len(preview.Images[0].Resolutions) <= 0 {
		return ""
	}

	return html.UnescapeString(preview.Images[0].Resolutions[0].URL)

}

func strNotEmpty(s string) bool {
	return len(s) > 0
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
