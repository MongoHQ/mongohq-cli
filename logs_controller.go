package main

import (
	"fmt"
	"strconv"
)

func (c *Controller) HistoricalLogs(deployment string) {
	historicalLogs, hostLength, err := c.Api.GetHistoricalLogs(deployment)

	if err != nil {
		fmt.Println("Error retrieving deployments: " + err.Error())
	} else {
		for _, log := range historicalLogs {
			fmt.Println(fmt.Sprintf("%-"+strconv.Itoa(hostLength+2)+"s%s", formatHostname(log.Host), log.Message))
		}
	}
}
