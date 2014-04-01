package api

import (
  "bytes"
  "net/http"
  "io/ioutil"
  "github.com/MongoHQ/mongohq-cli"
  "github.com/gorilla/websocket"
  "errors"
  "encoding/json"
  "strconv"
)

type OkResponse struct {
  Ok int
}

type ErrorResponse struct {
  Error string
}

var oauth_client_id = "6fb9368538ef061ed73be71cc291e65b"
var oauth_secret    = "028d31d8ca253cc3004b3ae4470c21bb23c3011e2fc8b442ad72f259be7879ce5c66bfda4ff26d5a0ba8d23369ef3355ef4579f6e7a977ba933dc1a37fd2880c"

func api_url(path string) (string) {
  return "https://dblayer-api.herokuapp.com" + path;
}

func gopher_url(path string) (string) {
  return "https://beta-api.mongohq.com/mongo" + path
}

func gopher_socket_url(path string, oauthToken string) (string) {
  return "wss://beta-api.mongohq.com/mongo" + path + "?token=Bearer%20" + oauthToken
}

func userAgent() string {
  return "MongoHQ CLI Version " + mongohq_cli.Version()
}

func sendRequest(request *http.Request, oauthToken string) ([]byte, error) {
  client := &http.Client{}

  request.Header.Add("Authorization", "Bearer " + oauthToken)
  request.Header.Add("User-Agent", userAgent())
  request.Header.Add("Content-Type", "application/json")
  response, err := client.Do(request)

  if err != nil {
    return nil, err
  }

  responseBody, _ := ioutil.ReadAll(response.Body)
  response.Body.Close()

  if string(responseBody) == "NOT FOUND" {
    return responseBody, errors.New("Object not found")
  } else if response.StatusCode >= 400 {
    var errorResponse ErrorResponse
    err := json.Unmarshal(responseBody, &errorResponse)

    if err != nil {
      return responseBody, err
    }
    return responseBody, errors.New("Response status " + response.Status + " with error " + errorResponse.Error)
  }

  return responseBody, nil
}

func rest_get(urlString, oauthToken string) ([]byte, error) {
  request, err := http.NewRequest("GET", urlString, nil)
  if err != nil {
    return nil, err
  }
  return sendRequest(request, oauthToken)
}

func rest_post(urlString string, data []byte, oauthToken string) ([]byte, error) {
  request, err := http.NewRequest("POST", urlString, bytes.NewReader(data))
  if err != nil {
    return nil, err
  }
  return sendRequest(request, oauthToken)
}

func rest_delete(urlString, oauthToken string) ([]byte, error) {
  request, err := http.NewRequest("DELETE", urlString, nil)
  if err != nil {
    return nil, err
  }
  return sendRequest(request, oauthToken)
}

func open_websocket(message SocketMessage, oauthToken string) (*websocket.Conn, error) {
  dialer := websocket.Dialer{}
  header := http.Header{}
  header.Add("User-Agent", userAgent())
  client, _, err := dialer.Dial(gopher_socket_url("/ws", oauthToken), header)
  if err != nil {
    return client, errors.New("Error initiating connection to websocket.")
  }
  jsonMessage, err := json.Marshal(message)
  if err != nil {
    return client, errors.New("Error marshalling websocket message.")
  }
  err = client.WriteMessage(websocket.TextMessage, jsonMessage)
  if err != nil {
    return client, errors.New("Error subscribing to websocket feed.")
  }
  return client, nil
}

var kb = 1024.0
var mb = kb * 1024.0
var gb = mb * 1024.0
var tb = gb * 1024.0

func prettySize(size float64) string {
  if size < kb {
    return strconv.FormatFloat(size, 'f', 0, 64) + "b"
  } else if size < mb {
    return strconv.FormatFloat(size / kb, 'f', 0, 64) + "kb"
  } else if size < gb {
    return strconv.FormatFloat(size / mb, 'f', 0, 64) + "mb"
  } else if size < tb {
    return strconv.FormatFloat(size / gb, 'f', 0, 64) + "gb"
  } else {
    return strconv.FormatFloat(size / tb, 'f', 0, 64) + "tb"
  }
}
