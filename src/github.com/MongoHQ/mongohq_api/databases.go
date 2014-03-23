package mongohq_api

import (
  "net/http"
  //"fmt"
  "io/ioutil"
  //"net/url"
  "encoding/json"
  //"errors"
)

type Database struct {
    Id   string
    Name string
}

func GetDatabases(oauthToken string) ([]Database, error) {
  client := &http.Client{}
  request, err := http.NewRequest("GET", rest_url_for("/databases"), nil)
  request.Header.Add("Authorization", "Bearer " + oauthToken)
  response, err := client.Do(request)

  if err != nil {
    return nil, err
  } else {
    responseBody, _ := ioutil.ReadAll(response.Body)
    response.Body.Close()

    var databasesSlice []Database
    _ = json.Unmarshal(responseBody, &databasesSlice)

    return databasesSlice, err
  }
}
