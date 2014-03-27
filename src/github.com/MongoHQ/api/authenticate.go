package api

import (
  "net/http"
  "net/url"
  "io/ioutil"
  "encoding/json"
  "errors"
  //"strings"
)

func Authenticate(username, password string) (string, error) {
  var oauthToken string
  var authenticationError error
  var jsonResponse map[string]interface{}

  data := make(url.Values)
  data.Set("username", username)
  data.Set("password", password)
  data.Set("grant_type", "password")
  data.Set("client_id", oauth_client_id)
  data.Set("client_secret", oauth_secret)

  response, err := http.PostForm(api_url("/login/oauth/access_token"), data)

  if err != nil {
    authenticationError = errors.New("Error authenticating against MongoHQ.")
  } else {
    responseBody, _ := ioutil.ReadAll(response.Body)
    _ = json.Unmarshal(responseBody, &jsonResponse)
    response.Body.Close()

    oauthToken = jsonResponse["access_token"].(string)
  }

  return oauthToken, authenticationError
}
