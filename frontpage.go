package main

import "github.com/labstack/echo/v4"

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