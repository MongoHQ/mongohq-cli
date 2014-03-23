package mongohq_cli

import (
  "fmt"
  "github.com/MongoHQ/mongohq_api"
)

func Databases() {
  databases, err := mongohq_api.GetDatabases(OauthToken)

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
