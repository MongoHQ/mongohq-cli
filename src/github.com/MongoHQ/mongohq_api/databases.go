package mongohq_api

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
  } else {
    var databasesSlice []Database
    _ = json.Unmarshal(body, &databasesSlice)

    return databasesSlice, err
  }
}
