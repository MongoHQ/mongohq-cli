package api

import (
  "encoding/json"
  "github.com/gorilla/websocket"
  "net/http"
  "os"
)

type Deployment struct {
    Id   string
    Current_primary string
    Version string
    Members []string
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

func GetDeployments(oauthToken string) ([]Deployment, error) {
  body, err := rest_get("/deployments", oauthToken)

  if err != nil {
    return nil, err
  } else {
    var deploymentsSlice []Deployment
    _ = json.Unmarshal(body, &deploymentsSlice)

    return deploymentsSlice, err
  }
}

func DeploymentMongostat(deployment_id string, database_name string, oauthToken string, outputFormatter func(string)) {
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
    messageType, p, err := client.ReadMessage()
    if err != nil {
      println("Error receiving message: " + err.Error())
      os.Exit(1)
    }
    println(string(messageType))
    println(string(p))
  }
}
