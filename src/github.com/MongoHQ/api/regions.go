package api

import (
  "encoding/json"
)

func GetRegions(oauthToken string) (map[string][]string, error) {
  body, err := rest_get(api_url("/deployment_regions"), oauthToken)

  if err != nil {
    return make(map[string][]string), err
  }
  var providersRegions map[string][]string
  err = json.Unmarshal(body, &providersRegions)
  return providersRegions, err
}
