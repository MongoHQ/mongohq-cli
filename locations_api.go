package main

import (
	"encoding/json"
)

func (api *Api) GetLocations() ([]string, error) {
	body, err := api.restGet(api.apiUrl("/locations"))

	if err != nil {
		return make([]string, 0), err
	}
	var providersLocations []string
	err = json.Unmarshal(body, &providersLocations)
	return providersLocations, err
}
