package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type AuthenticationArguments struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	GrantType string `json:"grant_type"`
	ClientId  string `json:"client_id"`
}

func (api Api) Authenticate(username, password, token string) (string, error) {
	var oauthToken string
	var authenticationError error
	var jsonResponse map[string]interface{}

	data, err := json.Marshal(AuthenticationArguments{Username: username, Password: password, GrantType: "password", ClientId: oauth_client_id})
	if err != nil {
		return oauthToken, errors.New("Error creating MongoHQ authentication request.")
	}

	client, err := buildHttpClient()
	if err != nil {
		return "", errors.New("Error building HTTPS transport process.")
	}
	request, err := http.NewRequest("POST", api.apiUrl("/oauth/token"), bytes.NewReader(data))

	if token != "" {
		request.Header.Add("X-Mongohq-Otp", token)
	}

	request.Header.Add("User-Agent", api.UserAgent)
	request.Header.Add("Content-Type", "application/json")

	response, err := client.Do(request)
	responseBody, _ := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", errors.New("Error authenticating against MongoHQ: " + err.Error())
	} else if response.StatusCode >= 400 {
		if response.Header.Get("X-Mongohq-Otp") == "required; sms" {
			return "", errors.New("2fa token required")
		} else if response.Header.Get("X-Mongohq-Otp") == "required; unconfigured" {
			return "", errors.New("Account requires 2fa authentication.  Go to https://app.mongohq.com to configure")
		} else {
			var errorResponse ErrorResponse
			err = json.Unmarshal(responseBody, &errorResponse)

			if err == nil {
				return "", errors.New(errorResponse.Error)
			}

			return "", errors.New("Error authenticating against MongoHQ.")
		}
	} else {
		_ = json.Unmarshal(responseBody, &jsonResponse)
		response.Body.Close()

		oauthToken = jsonResponse["access_token"].(string)
	}

	return oauthToken, authenticationError
}
