package controllers

import (
  "fmt"
  "github.com/MongoHQ/api"
)

func Databases() {
  databases, err := api.GetDatabases(OauthToken)

  if err != nil {
    fmt.Println("Error retrieving databases: " + err.Error())
  } else { 
    fmt.Println("=== My Databases")
    for _, database := range databases {
      fmt.Println(database.Name)
      fmt.Println("  status:        " + database.Status)
      fmt.Println("  plan:          " + database.Plan)
      fmt.Println("  deployment_id: " + database.Deployment_id)
    }
  }
}
