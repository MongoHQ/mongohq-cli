package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Api struct {
  OauthToken string
  UserAgent string
}

type OkResponse struct {
	Ok int
}

type ErrorResponse struct {
	Error string
}

var oauth_client_id = "5347253b33936795b1000001"

func (api *Api) apiUrl(path string) string {
	return "https://beta-api.mongohq.com/" + path
}

func (api *Api) gopherUrl(path string) string {
	return "https://beta-api.mongohq.com/mongo" + path
}

func (api *Api) gopherSocketUrl(path string) string {
	return "wss://beta-api.mongohq.com/mongo" + path + "?token=Bearer%20" + api.OauthToken
}

func (api *Api) sendRequest(request *http.Request) ([]byte, error) {
	client := &http.Client{}

  if api.OauthToken == "" {
    return nil, errors.New("Unknown oauth token.  Please run `mongohq logout`, then rerun your command.")
  }

	request.Header.Add("Authorization", "Bearer "+api.OauthToken)
	request.Header.Add("User-Agent", api.UserAgent)
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

func (api *Api) restGet(urlString string) ([]byte, error) {
	request, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		return nil, err
	}
	return api.sendRequest(request)
}

func (api *Api) restPost(urlString string, data []byte) ([]byte, error) {
	request, err := http.NewRequest("POST", urlString, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return api.sendRequest(request)
}

func (api *Api) restPatch(urlString string, data []byte) ([]byte, error) {
	request, err := http.NewRequest("PATCH", urlString, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return api.sendRequest(request)
}

func (api *Api) restDelete(urlString string) ([]byte, error) {
	request, err := http.NewRequest("DELETE", urlString, nil)
	if err != nil {
		return nil, err
	}
	return api.sendRequest(request)
}

func (api *Api) openWebsocket(message SocketMessage) (*websocket.Conn, error) {
	dialer := websocket.Dialer{}
	header := http.Header{}
	header.Add("User-Agent", api.UserAgent)
	client, _, err := dialer.Dial(api.gopherSocketUrl("/ws"), header)
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
