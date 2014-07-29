package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Controller struct {
	Api *Api
}

func prompt(text string) string {
	if replMode == true {
		line, _ := term.Prompt(text)
		return line
	}

	var response string
	var err error

	print(text + ": ")
	_, err = fmt.Scanln(&response)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return response
}

func formatHostname(host string) string {
	hostRegex := regexp.MustCompile(".(?:mongohq|mongolayer).com")
	host = hostRegex.ReplaceAllLiteralString(host, "")
	return host
}

func (c *Controller) pollNewDeployment(deployment Deployment) {
	var err error
	status := deployment.Status

	for status == "new" {
		fmt.Print(".")
		deployment, err = c.Api.GetDeployment(deployment.Name)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("\nError pulling deployment information.  For a manual update, please run:\n\n mongohq deployment:info --deployment " + deployment.Name)
			os.Exit(1)
		}
		status = deployment.Status
	}

	fmt.Print("\n")
	fmt.Println("Your database is ready. To add a user to your database, run:")
	fmt.Println("  mongohq users:create --deployment " + deployment.Name + " --database " + deployment.Databases[0].Name + " -u <username>")
	fmt.Println("")
	fmt.Println("To connect to your database, run:")
	fmt.Println("  mongo " + deployment.CurrentPrimary + "/" + deployment.Databases[0].Name + " -u <username>" + " -p")
	fmt.Println("")
	fmt.Println("Your applications should use the following URI to connect:")
	fmt.Println("  mongodb://<username>:<password>@" + strings.Join(deployment.Members, ",") + "/" + deployment.Databases[0].Name)
	fmt.Println("\nEnjoy!")
}
