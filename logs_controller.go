package main

import (
	"fmt"
	"strconv"
)

func (c *Controller) HistoricalLogs(deploymentSlug, search, exclude, regexp string) {
	historicalLogs, hostLength, err := c.Api.GetHistoricalLogs(deploymentSlug, search, exclude, regexp)

	if err != nil {
		fmt.Println("Error retrieving logs: " + err.Error())
	} else {
		for _, log := range historicalLogs {
			fmt.Println(fmt.Sprintf("%-"+strconv.Itoa(hostLength+2)+"s%s", formatHostname(log.Host), log.Message))
		}
	}
}
