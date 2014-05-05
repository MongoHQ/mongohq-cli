package main

import (
  //"fmt"
  "encoding/json"
  "io/ioutil"
  "os"
  "errors"
  "github.com/codegangsta/cli"
  "code.google.com/p/gopass"
  "fmt"
)

// Login is its on controller because it acts differently than others
type LoginController struct {
  Api Api
  OauthToken string
  Username string
}

var credentialPath = os.Getenv("HOME") + "/.mongohq"
var credentialFile = credentialPath + "/credentials"

var Email, OauthToken string

func (c *LoginController) login() (error) {
  fmt.Println("Enter your MongoHQ credentials.")
  username := prompt("Email")
  password, err := gopass.GetPass("Password (typing will be hidden): ")

  if err != nil {
    return errors.New("Error returning password.  We may not be compliant with your system yet.  Please send us a message telling us about your system to support@mongohq.com.")
  }

  oauthToken, err := c.Api.Authenticate(username, password, "")

  if err == nil {
    c.Api = Api{OauthToken: oauthToken}
  }

  return c.processAuthenticationResponse(username, password, oauthToken, err)
}

func (c *LoginController) processAuthenticationResponse(username, password, oauthToken string, err error) (error) { 
  if err != nil {
    if err.Error() == "2fa token required" {
      twoFactorToken := prompt("2fa token")
      oauthToken, err := c.Api.Authenticate(username, password, twoFactorToken)
      return c.processAuthenticationResponse(username, password, oauthToken, err)
    } else {
      return err
    }
  } else {
    err = c.storeCredentials(username, oauthToken)

    if err != nil {
      return err
    } else {
      fmt.Println("\nAuthentication complete.\n\n")
      c.OauthToken, c.Username = username, oauthToken
      return nil
    }
  }
}

func (c *LoginController) storeCredentials(username, oauth string) (error) {
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

func (c *LoginController) readCredentialFile() (jsonResponse map[string]interface{}, err error) {
  if _, err := os.Stat(credentialFile); os.IsNotExist(err) { // check if file exists
    return nil, errors.New("Credential file does not exist.")
  } else {
    jsonText, err := ioutil.ReadFile(credentialFile)
    _ = json.Unmarshal(jsonText, &jsonResponse)

    c.Api = Api{OauthToken: jsonResponse["oauth_token"].(string)}

    return jsonResponse, err
  }
}

func (c *LoginController) RequireAuth(*cli.Context) (err error) {
  for !c.verifyAuth() {}
  return err
}

func (c *LoginController) Logout() {
  os.Remove(credentialFile)
  fmt.Println("Logout success.")
}

func (c *LoginController) verifyAuth() (bool) {
  _, err := c.readCredentialFile()
  if err != nil {
    err := c.login()

    if err != nil {
      fmt.Println("\n"+err.Error()+"\n")
      return false
    } else {
      return true
    }
  } else {
     return true
  }
}
