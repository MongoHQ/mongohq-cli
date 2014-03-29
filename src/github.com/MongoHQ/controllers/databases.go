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

func CreateDatabase(deploymentId, databaseName string) {
  database, err := api.CreateDatabase(deploymentId, databaseName, OauthToken)

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

func DatabaseUsers(deploymentId, databaseName string) {
  databaseUsersSlice, err := api.GetDatabaseUsers(deploymentId, databaseName, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database: " + err.Error())
  } else {
    fmt.Println("== Users for database " + databaseName)
    for _, databaseUser := range databaseUsersSlice {
      fmt.Println("  " + databaseUser.Username)
    }
  }
}

func DatabaseCreateUser(deploymentId, databaseName, username string) {
  password := prompt("New user password")

  _, err := api.CreateDatabaseUser(deploymentId, databaseName, username, password, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database: " + err.Error())
  }
  fmt.Println("User " + username + " created.")
}

func DatabaseRemoveUser(deploymentId, databaseName, username string) {
  _, err := api.RemoveDatabaseUser(deploymentId, databaseName, username, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database: " + err.Error())
  }
  fmt.Println("User " + username + " removed.")
}
