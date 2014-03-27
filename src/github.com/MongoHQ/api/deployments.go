package api

import (
  "encoding/json"
  "github.com/gorilla/websocket"
  "net/http"
  "os"
  "strings"
)

type Deployment struct {
    Id   string
    CurrentPrimary string `json:"current_primary"`
    Version string
    Members []string
    AllowMultipleDatabases bool `json:"allow_multiple_deployments"`
}

type SocketMessage struct {
  Command string `json:"command"`
  Uuid string `json:"uuid"`
  Message Message `json:"message"`
}

type Message struct {
  DeploymentId string `json:"deployment_id"`
  DatabaseName string `json:"database_name"`
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
  body, err := rest_get("/deployments", oauthToken)

  if err != nil {
    return nil, err
  }
  var deploymentsSlice []Deployment
  err = json.Unmarshal(body, &deploymentsSlice)
  return deploymentsSlice, err
}

func GetDeployment(deploymentId string, oauthToken string) (Deployment, error) {
  body, err := rest_get("/deployments/" + deploymentId, oauthToken)

  if err != nil {
    return Deployment{}, err
  }
  var deployment Deployment
  err = json.Unmarshal(body, &deployment)
  return deployment, err
}

func DeploymentMongostat(deployment_id string, database_name string, oauthToken string, outputFormatter func([]map[string]MongoStat, error)) {
  message := SocketMessage {Command: "subscribe", Uuid: "12345", Message: Message{DeploymentId: deployment_id, DatabaseName: database_name, Type: "mongo.stats"}}

  dialer := websocket.Dialer{}
  header := http.Header{}
  header.Add("User-Agent", userAgent())
  client, _, err := dialer.Dial(socket_url_for("/ws", oauthToken), header)
  if err != nil {
    println("Error initiating connection to websocket: " + err.Error())
    os.Exit(1)
  }
  jsonMessage, err := json.Marshal(message)
  if err != nil {
    println("Error marshalling JSON: " + err.Error())
    os.Exit(1)
  }
  err = client.WriteMessage(websocket.TextMessage, jsonMessage)
  if err != nil {
    println("Error subscribing to feed: " + err.Error())
    os.Exit(1)
  }

  for {
    _, msg, err := client.ReadMessage()
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
