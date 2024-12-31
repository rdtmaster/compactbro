package main

import (
	"fmt"
	"html"
	"html/template"
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
	"github.com/rdtmaster/go-reddit/v4/reddit"
)

const emptyHTML = template.HTML("")
var (
	config CompactConfig
	server *echo.Echo
	tpls map[string]*template.Template
	client *reddit.Client
)

// Shut the server down
func shutdown(c echo.Context) error {
	go func() {
		time.Sleep(500 * time.Millisecond)
		server.Close()
		os.Exit(0)
	}()
	return c.HTML(200, `<html><head><title>Goodbye</title><body bgcolor="#000123" text="#cdedfe"><h1 style="font-family: tahoma;right: 50%;bottom: 50%;transform: translate(50%,50%);position: absolute">Stopping server...</h1></body></html>`)
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

	amber.FuncMap["linkFromContext"] = linkFromContext
	amber.FuncMap["cleanCommentID"] = cleanCommentID
	amber.FuncMap["cleanLink"] = cleanLink
	amber.FuncMap["isPostID"] = isPostID
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
	server.HideBanner = true
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
	server.GET("/message/", messageInbox)
	server.GET("/message/:page/", messageInbox)
	server.GET("/pt/message/:page/", messageInboxPT)
	server.GET("/pt/", frontpagePT)
	server.GET("/r/:sub", subDefault)
	server.GET("/r/:sub/", subDefault)
	server.GET("/r/:sub/:sorting", subSorted)
	server.GET("/r/:sub/:sorting/", subSorted)

	server.GET("/pt/r/:sub/", subPT)
	server.GET("/r/:sub/comments/:id/", submission)
	server.GET("/r/:sub/comments/:id/:permalink/", submission)
	server.GET("/user/:sub/comments/:id/:permalink/", submission)
	server.GET("/u/:username/", overview)
	server.GET("/u/:username", overview)
	server.GET("/u/:username/:page/", overview)
	server.GET("/pt/u/:username/:page/", overviewPT)
	server.POST("/edit/", editThing)
	server.POST("/comment/", submitComment)
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