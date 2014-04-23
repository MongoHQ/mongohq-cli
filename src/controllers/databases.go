package controllers

import (
  "fmt"
  "api"
  "os"
)

func Databases() {
  databases, err := api.GetDatabases(OauthToken)

  if err != nil {
    fmt.Println("Error retrieving databases: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("=== My Databases")
  for _, database := range databases {
    fmt.Println(database.Name)
  }
}

func Database(name string) {
  database, err := api.GetDatabase(name, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("=== " + database.Name)
  fmt.Println(" id:            " + database.Id)
  fmt.Println(" name:          " + database.Name)
  fmt.Println(" plan:          " + database.Plan)
  fmt.Println(" status:        " + database.Status)
  fmt.Println(" deployment_id: " + database.Deployment_id)
}

func RemoveDatabase(databaseName string, force bool) {
  if !force {
    confirmDatabaseName := prompt("To confirm, type the name of the database to be deleted")

    if databaseName != confirmDatabaseName {
      fmt.Println("Confirmation of database name is incorrect.")
      os.Exit(1)
    }
  }

  err := api.RemoveDatabase(databaseName, OauthToken)

  if err != nil {
    fmt.Println("Error removing database: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("Removed database named: " + databaseName)
}

func CreateDatabase(deploymentId, databaseName string) {
  database, err := api.CreateDatabase(deploymentId, databaseName, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("=== " + database.Name)
  fmt.Println(" id:            " + database.Id)
  fmt.Println(" name:          " + database.Name)
  fmt.Println(" plan:          " + database.Plan)
  fmt.Println(" status:        " + database.Status)
  fmt.Println(" deployment_id: " + database.Deployment_id)
}

func DatabaseUsers(deploymentId, databaseName string) {
  databaseUsersSlice, err := api.GetDatabaseUsers(deploymentId, databaseName, OauthToken)

  if err != nil {
    fmt.Println("Error retrieiving database users: " + err.Error())
    os.Exit(1)
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
    fmt.Println("Error creating database user: " + err.Error())
    os.Exit(1)
  }
  fmt.Println("User " + username + " created.")
}

func DatabaseRemoveUser(deploymentId, databaseName, username string) {
  _, err := api.RemoveDatabaseUser(deploymentId, databaseName, username, OauthToken)

  if err != nil {
    fmt.Println("Error removing database users: " + err.Error())
    os.Exit(1)
  }
  fmt.Println("User " + username + " removed.")
}
