package main

import (
  "encoding/json"
)

func (api *Api) GetRegions() (map[string][]string, error) {
  body, err := api.restGet(api.apiUrl("/deployment_regions"))

  if err != nil {
    return make(map[string][]string), err
  }
  var providersRegions map[string][]string
  err = json.Unmarshal(body, &providersRegions)
  return providersRegions, err
}
