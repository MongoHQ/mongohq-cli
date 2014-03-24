package mongohq_cli

import (
  "fmt"
  "github.com/MongoHQ/mongohq_api"
  "strings"
)

func Deployments() {
  deployments, err := mongohq_api.GetDeployments(OauthToken)

  if err != nil {
    fmt.Println("Error retrieving deployments: " + err.Error())
  } else {
    fmt.Println("=== My Deployments")
    for _, deployment := range deployments {
      fmt.Println(deployment.Current_primary)
      fmt.Println("  Id:              " + deployment.Id)
      fmt.Println("  MongoDB Version: " + deployment.Version)
      fmt.Println("  Members:         " + strings.Join(deployment.Members, ","))
    }
  }
}

func MongoStat(deployment_id, database_name string) {
  outputFormatter := func(msg string) {
  }

  mongohq_api.DeploymentMongostat(deployment_id, database_name, OauthToken, outputFormatter)
}
