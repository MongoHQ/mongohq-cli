package mongohq_api

import (
  "encoding/json"
)

type Deployment struct {
    Id   string
    Current_primary string
    Version string
    Members []string
}

func GetDeployments(oauthToken string) ([]Deployment, error) {
  body, err := rest_get("/deployments", oauthToken)

  if err != nil {
    return nil, err
  } else {
    var deploymentsSlice []Deployment
    _ = json.Unmarshal(body, &deploymentsSlice)

    return deploymentsSlice, err
  }
}
