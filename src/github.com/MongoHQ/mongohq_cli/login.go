package mongohq_cli

import (
  "github.com/MongoHQ/mongohq_api"
  "fmt"
  "encoding/json"
  "io/ioutil"
  "os"
  "errors"
)

var credential_path = os.Getenv("HOME") + "/.mongohq"
var credential_file = credential_path + "/credentials"

func Login() {
  username := prompt("Username")
  password := prompt("Password")

  //println("2fa access required.  Enter your 2fa token: ")
  //var 2fa string
  //_, err = fmt.Scanln(&2fa)
  //if err != nil {
    //fmt.Println("Error: ", err)
  //}
  oauth_token, err := mongohq_api.Authenticate(username, password)

  if err != nil {
    fmt.Println("Error authenticating")
    //fmt.Println("Error: ", err.Error())
  } else {
    err = storeCredentials(username, oauth_token)

    if err != nil {
      fmt.Println(err)
    } else {
      fmt.Println("Authentication complete.")
    }
  }
}

func storeCredentials(username, oauth string) (error) {
  credentials := make(map[string]interface{})
  credentials["email"] = username
  credentials["oauth_token"] = oauth

  jsonText, _ := json.Marshal(credentials)

  err := os.MkdirAll(credential_path, 0700)

  if err != nil {
    return errors.New("Error creating directory " + credential_path)
  }

  err = ioutil.WriteFile(credential_file, jsonText, 0500)

  if err != nil {
    err = errors.New("Error writing credential_file to " + credential_file)
  }

  return err
}

func GetCredentials() (jsonResponse map[string]interface{}, err error) {
  jsonText, err := ioutil.ReadFile(credential_file)
  _ = json.Unmarshal(jsonText, &jsonResponse)
  return jsonResponse, err
}
