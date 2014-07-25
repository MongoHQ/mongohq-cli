package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func (c *Controller) HistoricalLogs(deploymentSlug, search, exclude, regexp string) {
	logLimit := 200

	var first time.Time
	var last time.Time

	var command string

	for command != "exit" && command != "e" {
		var historicalLogs []HistoricalLog
		var hostLength int
		var err error

		if command == "n" || command == "next" {
			historicalLogs, hostLength, err = c.Api.GetHistoricalLogs(deploymentSlug, search, exclude, regexp, logLimit, &last, nil)
		} else if command == "p" || command == "previous" {
			historicalLogs, hostLength, err = c.Api.GetHistoricalLogs(deploymentSlug, search, exclude, regexp, logLimit, nil, &first)
		} else if command == "" {
			historicalLogs, hostLength, err = c.Api.GetHistoricalLogs(deploymentSlug, search, exclude, regexp, logLimit, nil, nil)
		} else {
			fmt.Print(command + " is an unknown option.  Please type (n)ext, (p)revious, or (e)xit.")
			continue
		}

		if err != nil {
			fmt.Println("Error retrieving logs: " + err.Error())
			os.Exit(1)
		} else {
			if len(historicalLogs) == 0 {
				fmt.Println("No logs matching query.")
			} else {
				first, last = renderLogs(historicalLogs, hostLength)
			}
		}

		command = prompt("(n)ext (p)revious (e)xit >")
	}
}

func renderLogs(historicalLogs []HistoricalLog, hostLength int) (time.Time, time.Time) {
	var last HistoricalLog

	first := historicalLogs[0]
	for _, log := range historicalLogs {
		last = log
		fmt.Println(fmt.Sprintf("%-"+strconv.Itoa(hostLength+2)+"s%s", formatHostname(log.Host), log.Message))
	}

	return first.Timestamp.Add(time.Millisecond * -1), last.Timestamp.Add(time.Millisecond)
}
