package api

import (
  "encoding/json"
)

type Database struct {
    Id   string `json:"id"`
    Name string `json:"name"`
    Status string `json:"status"`
    Plan string `json:"plan"`
    Deployment_id string `json:"deployment_id"`
}

type DatabaseUser struct {
    Username string `json:"user"`
    PasswordHash string `json:"pwd"`
    ReadOnly bool
}

func GetDatabases(oauthToken string) ([]Database, error) {
  body, err := rest_get(api_url("/databases"), oauthToken)

  if err != nil {
    return nil, err
  }
  var databasesSlice []Database
  err = json.Unmarshal(body, &databasesSlice)
  return databasesSlice, err
}

func GetDatabase(name string, oauthToken string) (Database, error) {
  body, err := rest_get(api_url("/databases/" + name), oauthToken)

  if err != nil {
    return Database{}, err
  }
  var database Database
  err = json.Unmarshal(body, &database)
  return database, err
}

func CreateDatabase(deploymentId, databaseName, oauthToken string) (Database, error) {
  type DatabaseCreate struct {
    DeploymentId string `json:"deployment_id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
  }

  databaseCreate := DatabaseCreate{Name: databaseName, Slug: "general:on_existing", DeploymentId: deploymentId}
  data, err := json.Marshal(databaseCreate)
  if err != nil {
    return Database{}, err
  }

  body, err := rest_post(api_url("/databases"), data, oauthToken)

  if err != nil {
    return Database{}, err
  }
  var database Database
  err = json.Unmarshal(body, &database)
  return database, err
}

func RemoveDatabase(databaseName, oauthToken string) error {
  _, err := rest_delete(api_url("/databases/" + databaseName), oauthToken)
  return err
}

func GetDatabaseUsers(deployment_id, database_name, oauthToken string) ([]DatabaseUser, error) {
  body, err := rest_get(gopher_url("/" + deployment_id + "/" + database_name + "/users"), oauthToken)
  if err != nil {
    return make([]DatabaseUser, 0), err
  }
  var databaseUsersSlice []DatabaseUser
  err = json.Unmarshal(body, &databaseUsersSlice)
  return databaseUsersSlice, err
}

func CreateDatabaseUser(deploymentId, databaseName, username, password, oauthToken string) (OkResponse, error) {
  type UserCreate struct {
    Username string `json:"username"`
    Password string `json:"password"`
    ReadOnly bool `json:"readOnly"`
  }

  userCreate := UserCreate{Username: username, Password: password, ReadOnly: false}
  data, err  := json.Marshal(userCreate)
  if err != nil {
    return OkResponse{}, err
  }

  body, err := rest_post(gopher_url("/" + deploymentId + "/" + databaseName + "/users"), data, oauthToken)

  if err != nil {
    return OkResponse{}, err
  }
  var okResponse OkResponse
  err = json.Unmarshal(body, &okResponse)
  return okResponse, err
}

func RemoveDatabaseUser(deploymentId, databaseName, username, oauthToken string) (OkResponse, error) {
  body, err := rest_delete(gopher_url("/" + deploymentId + "/" + databaseName + "/users/" + username), oauthToken)

  if err != nil {
    return OkResponse{}, err
  }
  var okResponse OkResponse
  err = json.Unmarshal(body, &okResponse)
  return okResponse, err
}
