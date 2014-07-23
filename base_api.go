package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type Api struct {
	OauthToken string
	UserAgent  string
	Config     *Config
}

type Hateos struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type OkResponse struct {
	Ok int
}

type ErrorResponse struct {
	Error string
}

func (api *Api) apiUrl(path string) string {
	return "https://api.mongohq.com" + path
}

func (api *Api) gopherSocketUrl(path string) string {
	return "wss://api.mongohq.com/mongo" + path + "?token=Bearer%20" + api.OauthToken
}

func decodePem(certInput string) tls.Certificate {
	var cert tls.Certificate
	certPEMBlock := []byte(certInput)
	var certDERBlock *pem.Block
	for {
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		if certDERBlock != nil && certDERBlock.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDERBlock.Bytes)
		}
	}
	return cert
}

func (api *Api) sendRequest(request *http.Request) ([]byte, error) {
	client, err := buildHttpClient()

	if err != nil {
		return nil, errors.New("Error building HTTPS transport process.")
	}

	if api.OauthToken == "" {
		return nil, errors.New("Unknown oauth token.  Please run `mongohq logout`, then rerun your command.")
	}

	request.Header.Add("Authorization", "Bearer "+api.OauthToken)
	request.Header.Add("User-Agent", api.UserAgent)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Version", "2014-06")
	response, err := client.Do(request)

	if err != nil {
		noConnection := regexp.MustCompile("no such host")
		if noConnection.Match([]byte(err.Error())) {
			fmt.Println("We couldn't find the MongoHQ host.  Typically, this means your internet connections has gone AWOL.")
			os.Exit(1)
		} else {
			return nil, err
		}
	}

	responseBody, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()

	if response.StatusCode >= 400 { // test for {error: "message"} type responses.
		var errorResponse ErrorResponse
		err = json.Unmarshal(responseBody, &errorResponse)

		if err == nil {
			return responseBody, errors.New(errorResponse.Error)
		}
	}

	if string(responseBody) == "NOT FOUND" {
		return responseBody, errors.New("Object not found")
	} else if response.StatusCode == 500 {
		return responseBody, errors.New("MongoHQ service returned an error. Check your parameters and try again, or check our status page: https://status.mongohq.com.")
	} else if response.StatusCode >= 401 {
		return responseBody, errors.New("Could not access the requested object.  Double check the arguments, or run `mongohq logout` and re-run the prior command.")
	} else if response.StatusCode >= 400 {
		var errorResponse ErrorResponse
		err := json.Unmarshal(responseBody, &errorResponse)

		if err != nil {
			return responseBody, err
		}
		return responseBody, errors.New("Response status " + response.Status + " with error " + errorResponse.Error)
	} else if response.Header.Get("X-User-Agent-Deprecated") == "true" {
		return responseBody, errors.New(response.Header.Get("X-User-Agent-Deprecation-Message"))
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
		return client, errors.New("could not initiate connection to websocket.")
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

func buildHttpClient() (http.Client, error) {
	certChain := decodePem(chain)
	conf := tls.Config{}
	conf.RootCAs = x509.NewCertPool()

	for _, cert := range certChain.Certificate {
		x509Cert, err := x509.ParseCertificate(cert)
		if err != nil {
			return http.Client{}, err
		}
		conf.RootCAs.AddCert(x509Cert)
	}
	conf.BuildNameToCertificate()

	tr := http.Transport{TLSClientConfig: &conf}

	return http.Client{Transport: &tr}, nil
}
