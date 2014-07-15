package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var configFile = configPath + "/defaults"

type Config struct {
	AccountSlug    string `json:"account-slug"`
	DeploymentSlug string `json:"deployment-slug"`
}

func getConfig() Config {
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return Config{}
	} else {
		var config Config

		jsonText, err := ioutil.ReadFile(configFile)

		if err != nil {
			return Config{}
		}

		_ = json.Unmarshal(jsonText, &config)

		return config
	}
}

func (d *Config) Save() error {
	jsonText, _ := json.Marshal(d)
	return ioutil.WriteFile(configFile, jsonText, 0600)
}

func (d *Config) requireAccount() error {
	if d.AccountSlug == "<string>" {
		return errors.New("Default account is required.  Please run `mongohq accounts` and `mongohq config:account -a <account-slug>` to set a default acount.")
	}
	return nil
}
