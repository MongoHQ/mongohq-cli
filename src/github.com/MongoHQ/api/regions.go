package api

import (
  "encoding/json"
)

func GetRegions(oauthToken string) ([]string, error) {
  body, err := rest_get(api_url("/deployment_regions"), oauthToken)

  if err != nil {
    return make([]string, 0), err
  }
  var regionsSlice []string
  err = json.Unmarshal(body, &regionsSlice)
  return regionsSlice, err
}
