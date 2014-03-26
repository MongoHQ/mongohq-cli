package controllers

import (
  "fmt"
  "github.com/MongoHQ/api"
  "strings"
)

func Deployments() {
  deployments, err := api.GetDeployments(OauthToken)

  if err != nil {
    fmt.Println("Error retrieving deployments: " + err.Error())
  } else {
    fmt.Println("=== My Deployments")
    for _, deployment := range deployments {
      fmt.Println(deployment.CurrentPrimary + " :: " + deployment.Id)
    }
  }
}

func Deployment(deploymentId string) {
  deployment, err := api.GetDeployment(deploymentId, OauthToken)

  if err != nil {
    fmt.Println("Error retrieving deployments: " + err.Error())
  } else {
    fmt.Println("=== " + deployment.Id)
    fmt.Println("  current primary:     " + deployment.CurrentPrimary)
    fmt.Println("  members:             " + strings.Join(deployment.Members, ","))
    fmt.Println("  version:             " + deployment.Version)

    if deployment.AllowMultipleDatabases { 
      fmt.Println("  multiple databases?: true")
    }
  }
}

func MongoStat(deployment_id, database_name string) {
  outputFormatter := func(msg string) {
  }

  api.DeploymentMongostat(deployment_id, database_name, OauthToken, outputFormatter)
}
