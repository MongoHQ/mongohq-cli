package controllers

import (
  //"fmt"
  "github.com/MongoHQ/api"
  "encoding/json"
  "io/ioutil"
  "os"
  "errors"
  "github.com/codegangsta/cli"
)

var credentialPath = os.Getenv("HOME") + "/.mongohq"
var credentialFile = credentialPath + "/credentials"

var Email, OauthToken string

func login() (string, string, error) {
  username := prompt("Username")
  password := prompt("Password")

  //println("2fa access required.  Enter your 2fa token: ")
  //var 2fa string
  //_, err = fmt.Scanln(&2fa)
  //if err != nil {
    //fmt.Println("Error: ", err)
  //}
  oauthToken, err := api.Authenticate(username, password)

  if err != nil {
    return username, "", errors.New("Error authenticating given username / password")
  } else {
    err = storeCredentials(username, oauthToken)

    if err != nil {
      return username, oauthToken, err
    } else {
      return username, oauthToken, nil
    }
  }
}

func storeCredentials(username, oauth string) (error) {
  credentials := make(map[string]interface{})
  credentials["email"] = username
  credentials["oauth_token"] = oauth

  jsonText, _ := json.Marshal(credentials)

  err := os.MkdirAll(credentialPath, 0700)

  if err != nil {
    return errors.New("Error creating directory " + credentialPath)
  }

  err = ioutil.WriteFile(credentialFile, jsonText, 0500)

  if err != nil {
    err = errors.New("Error writing credentials to " + credentialFile)
  }

  return err
}

func readCredentialFile() (jsonResponse map[string]interface{}, err error) {
  if _, err := os.Stat(credentialFile); os.IsNotExist(err) { // check if file exists
    return nil, errors.New("Credential file does not exist.")
  } else {
    jsonText, err := ioutil.ReadFile(credentialFile)
    _ = json.Unmarshal(jsonText, &jsonResponse)

    return jsonResponse, err
  }
}

func RequireAuth(*cli.Context) (err error) {
  for !verifyAuth() {}
  return err
}

func Logout() {
  os.Remove(credentialFile)
}

func verifyAuth() (bool) {
  userMap, err := readCredentialFile()
  if err != nil {
    username, oauthToken, err := login()

    if err != nil {
      return false
    } else {
      Email = username
      OauthToken = oauthToken
      return true
    }
  } else {
     Email = userMap["email"].(string)
     OauthToken = userMap["oauth_token"].(string)
     return true
  }
}
