package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"mongohq-cli"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
)

type OkResponse struct {
	Ok int
}

type ErrorResponse struct {
	Error string
}

var oauth_client_id = "5347253b33936795b1000001"

func api_url(path string) string {
	return "https://beta-api.mongohq.com/" + path
}

func gopher_url(path string) string {
	return "https://beta-api.mongohq.com/mongo" + path
}

func gopher_socket_url(path string, oauthToken string) string {
	return "wss://beta-api.mongohq.com/mongo" + path + "?token=Bearer%20" + oauthToken
}

func userAgent() string {
	return "MongoHQ CLI Version " + mongohq_cli.Version()
}

func sendRequest(request *http.Request, oauthToken string) ([]byte, error) {
	client := &http.Client{}

	request.Header.Add("Authorization", "Bearer "+oauthToken)
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

func rest_patch(urlString string, data []byte, oauthToken string) ([]byte, error) {
	request, err := http.NewRequest("PATCH", urlString, bytes.NewReader(data))
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

func includeSignificantDigits(size float64) string {
	if size < 10 {
		return strconv.FormatFloat(size, 'f', 2, 64)
	} else if size < 100 {
		return strconv.FormatFloat(size, 'f', 1, 64)
	} else {
		return strconv.FormatFloat(size, 'f', 0, 64)
	}
}

func prettySize(size float64) string {
	if size < kb {
		return includeSignificantDigits(size) + "b"
	} else if size < mb {
		return includeSignificantDigits(size/kb) + "k"
	} else if size < gb {
		return includeSignificantDigits(size/mb) + "m"
	} else if size < tb {
		return includeSignificantDigits(size/gb) + "g"
	} else {
		return includeSignificantDigits(size/tb) + "t"
	}
}
