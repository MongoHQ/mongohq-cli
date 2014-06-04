package main

import (
  "encoding/json"
  "io/ioutil"
  "os"
)

var defaultsFile = configPath + "/defaults"

type Defaults struct {
  Account string `json:"account-slug"`
  Deployment string `json:"deployment"`
}

func getDefaults() (Defaults) {
  if _, err := os.Stat(defaultsFile); os.IsNotExist(err) {
    return Defaults{}
  } else {
    var defaults Defaults

    jsonText, err := ioutil.ReadFile(defaultsFile)

    if err != nil {
      return Defaults{}
    }

    _ = json.Unmarshal(jsonText, &defaults)

    return defaults
  }
}

func (d *Defaults) Save () (error) {
  jsonText, _ := json.Marshal(d)
  return ioutil.WriteFile(defaultsFile, jsonText, 0600)
}

func SetDefaultAccount(slug string) (error) {
  defaults := getDefaults()
  defaults.Account = slug
  return defaults.Save()
}
