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
    }
  }
}

func Database(name string) {
  database, err := api.GetDatabase(name, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database: " + err.Error())
  } else {
    fmt.Println("=== " + database.Name)
    fmt.Println(" id:            " + database.Id)
    fmt.Println(" name:          " + database.Name)
    fmt.Println(" plan:          " + database.Plan)
    fmt.Println(" status:        " + database.Status)
    fmt.Println(" deployment_id: " + database.Deployment_id)
  }
}
