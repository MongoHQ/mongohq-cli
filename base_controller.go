package main

import (
	"fmt"
	"regexp"
  "os"
  "strings"
)

type Controller struct {
  Api Api
}

func prompt(text string) string {
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

func (c *Controller) pollNewDeployment(databaseRecord Database) {
  var database Database
  var err error
  status, database := databaseRecord.Status, databaseRecord

  for status == "new" {
    fmt.Print(".")
    database, err = c.Api.GetDatabase(database.Name)
    if err != nil {
      fmt.Println(err.Error())
      fmt.Println("\nError pulling database information.  For a manual update, please run:\n\n mongohq databases:info --database " + database.Name)
      os.Exit(1)
    }
    status = database.Status
  }

  deployment, err := c.Api.GetDeployment(database.Deployment_id)

  if err != nil {
    fmt.Println("\nError pulling new deployment information.  For a manual update, please run:\n\n mongohq databases:info --database " + database.Name)
    os.Exit(1)
  }

  fmt.Print("\n")
  fmt.Println("Your database is ready. To add a user to your database, run:")
  fmt.Println("  mongohq users:create --deployment " + database.Deployment_id + " --database " + database.Name + " -u <username>")
  fmt.Println("")
  fmt.Println("To connect to your database, run:")
  fmt.Println("  mongo " + deployment.CurrentPrimary + "/" + database.Name + " -u <username>" + " -p")
  fmt.Println("")
  fmt.Println("Your applications should use the following URI to connect:")
  fmt.Println("  mongodb://<username>:<password>@" + strings.Join(deployment.Members, ",") + "/" + database.Name)
  fmt.Println("\nEnjoy!")
}
