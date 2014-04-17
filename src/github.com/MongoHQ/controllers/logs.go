package controllers

import (
	"fmt"
	"github.com/MongoHQ/api"
)

func HistoricalLogs(deploymentId string) {
	historicalLogs, err := api.GetHistoricalLogs(deploymentId, OauthToken)

	if err != nil {
		fmt.Println("Error retrieving deployments: " + err.Error())
	} else {
		for _, log := range historicalLogs {
			fmt.Println(formatHostname(log.Host) + "  " + log.Message)
		}
		/* fmt.Println("=== " + deployment.Id)
		fmt.Println("  current primary:     " + deployment.CurrentPrimary)
		fmt.Println("  members:             " + strings.Join(deployment.Members, ","))
		fmt.Println("  version:             " + deployment.Version)

		fmt.Println("  == Databases")
		for _, database := range deployment.Databases {
			fmt.Println("    " + database.Name)
		} */
	}
}
