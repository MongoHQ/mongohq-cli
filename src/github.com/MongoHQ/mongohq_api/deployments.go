package mongohq_api

import (
  "net/http"
  //"fmt"
  "io/ioutil"
  //"net/url"
  "encoding/json"
  //"errors"
)

type Deployment struct {
    Id   string
    Current_primary string
    Version string
    Members []string
}

func GetDeployments(oauthToken string) ([]Deployment, error) {
  client := &http.Client{}
  request, err := http.NewRequest("GET", rest_url_for("/deployments"), nil)
  request.Header.Add("Authorization", "Bearer " + oauthToken)
  response, err := client.Do(request)

  if err != nil {
    return nil, err
  } else {
    responseBody, _ := ioutil.ReadAll(response.Body)
    response.Body.Close()

    var deploymentsSlice []Deployment
    _ = json.Unmarshal(responseBody, &deploymentsSlice)

    return deploymentsSlice, err
  }
}
