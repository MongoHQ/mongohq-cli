package api

import (
  "encoding/json"
  "strings"
)

type Deployment struct {
    Id   string
    CurrentPrimary string `json:"current_primary"`
    Version string
    Members []string
    AllowMultipleDatabases bool `json:"allow_multiple_deployments"`
    Databases []Database
}

type SocketMessage struct {
  Command string `json:"command"`
  Uuid string `json:"uuid"`
  Message Message `json:"message"`
}

type Message struct {
  DeploymentId string `json:"deployment_id"`
  Type string `json:"type"`
}

type MongoStatMessage struct {
  Type string
  Ts string
  Error string
  Message []map[string]MongoStat
}

type MongoStat struct {
  Version string
  Inserts string
  RawInserts int `json:"raw_inserts"`
  Query string
  RawQuery int `json:"raw_query"`
  Update string
  RawUpdate int `json:"raw_update"`
  Delete string
  RawDelete int `json:"raw_delete"`
  Getmore string
  RawGetmore int `json:"raw_getmore"`
  Command string
  RawCommand int `json:"raw_command"`
  Flushes int
  Mapped int
  Vsize int
  Res int
  Faults int
  Locked string
  IdxMiss int `json:"idx_miss"`
  Qr int
  Qw int
  Ar int
  Aw int
  NetIn float64 `json:"net_in"`
  NetOut float64 `json:"net_out"`
  Conn int
  Set string
  Repl string
}

func GetDeployments(oauthToken string) ([]Deployment, error) {
  body, err := rest_get(api_url("/deployments"), oauthToken)

  if err != nil {
    return nil, err
  }
  var deploymentsSlice []Deployment
  err = json.Unmarshal(body, &deploymentsSlice)
  return deploymentsSlice, err
}

func GetDeployment(deploymentId string, oauthToken string) (Deployment, error) {
  body, err := rest_get(api_url("/deployments/" + deploymentId) + "?embed=databases", oauthToken)

  if err != nil {
    return Deployment{}, err
  }
  var deployment Deployment
  err = json.Unmarshal(body, &deployment)
  return deployment, err
}

func CreateDeployment(deploymentName, databaseName, region, oauthToken string) (Database, error) {
  type DatabaseCreateOptions struct {
    Region string `json:"region"`
  }

  type DatabaseCreate struct {
    Name string `json:"name"`
    Slug string `json:"slug"`
    Options DatabaseCreateOptions
  }

  databaseCreate := DatabaseCreate{Name: databaseName, Slug: "mongohq:elastic", Options: DatabaseCreateOptions{Region: region}}
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

func DeploymentMongostat(deployment_id string, oauthToken string, outputFormatter func([]map[string]MongoStat, error)) error {
  message := SocketMessage {Command: "subscribe", Uuid: "12345", Message: Message{DeploymentId: deployment_id, Type: "mongo.stats"}}
  socket, err := open_websocket(message, oauthToken)
  if err != nil {
    return err
  }

  for {
    _, msg, err := socket.ReadMessage()
    if err != nil {
      outputFormatter(make([]map[string]MongoStat, 0), err)
    }

    // catch the first success response
    if strings.Index(string(msg), "successful") > -1 {
      continue
    }

    // null is bad news for Go, and gopher has an outstanding issue with
    // the first response: https://github.com/MongoHQ/gopher/issues/14
    if strings.Index(string(msg), "null") > -1 {
      continue
    }

    mongoStatMessage := MongoStatMessage{}
    err = json.Unmarshal(msg, &mongoStatMessage)
    outputFormatter(mongoStatMessage.Message, err)
  }
}

func DeploymentOplog(deployment_id string, oauthToken string, outputFormatter func(string, error)) error {
  message := SocketMessage {Command: "subscribe", Uuid: "12345", Message: Message{DeploymentId: deployment_id, Type: "mongo.oplog"}}
  socket, err := open_websocket(message, oauthToken)
  if err != nil {
    return err
  }

  for {
    _, msg, err := socket.ReadMessage()
    outputFormatter(string(msg), err)
  }
}
