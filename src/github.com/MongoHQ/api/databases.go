package api

import (
  "encoding/json"
)

type Database struct {
    Id   string
    Name string
    Status string
    Plan string
    Deployment_id string
}

func GetDatabases(oauthToken string) ([]Database, error) {
  body, err := rest_get("/databases", oauthToken)

  if err != nil {
    return nil, err
  }
  var databasesSlice []Database
  err = json.Unmarshal(body, &databasesSlice)
  return databasesSlice, err
}

func GetDatabase(name string, oauthToken string) (Database, error) {
  body, err := rest_get("/databases/" + name, oauthToken)

  if err != nil {
    return Database{}, err
  }
  var database Database
  err = json.Unmarshal(body, &database)
  return database, err
}
