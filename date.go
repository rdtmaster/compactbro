package main

import (
	"time"
	"fmt"
	"math"
	"github.com/rdtmaster/go-reddit/v4/reddit"
)

// Credit: https://www.socketloop.com/tutorials/golang-get-time-duration-in-year-month-week-or-day


func roundTime(input float64) int {
	var result float64
	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	i, _ := math.Modf(result)

	return int(i)
}

func calcrDate(diff time.Duration) (dVal int, dName string) {
	dName = ""
	dVal = 0
	secs := diff.Seconds()
	years := roundTime(secs / 31207680)
	if years > 0 {
		dVal = years
		dName = "year"
		return
	}

	months := roundTime(secs / 2600640)
	if months > 0 {
		dVal = months
		dName = "month"
		return
	}
	weeks := roundTime(secs / 604800)
	if weeks > 0 {
		dVal = weeks
		dName = "week"
		return
	}
	days := roundTime(secs / 86400)
	if days > 0 {
		dVal = days
		dName = "day"
		return
	}
	hours := roundTime(diff.Hours())
	if hours > 0 {
		dVal = hours
		dName = "hour"
		return
	}
	minutes := roundTime(diff.Minutes())
	if minutes > 0 {
		dVal = minutes
		dName = "minute"
		return
	}
	return
}
func dateAgo(t *reddit.Timestamp) string {
	dVal, dName := calcrDate(time.Since(t.Time))
	switch dVal {
	case 0:
		return "just now"
	case 1:
		return fmt.Sprintf("%d %s ago", dVal, dName)
	default:
		return fmt.Sprintf("%d %ss ago", dVal, dName)
	}
}
