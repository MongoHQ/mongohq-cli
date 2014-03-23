package mongohq_api

import (
  "net/http"
  "fmt"
  //"net/url"
  //"encoding/json"
  //"errors"
)

type Database struct {
    Id   string
    Name string
}

func Databases() ([]Database, error) {
  client := &http.Client{
    //CheckRedirect: redirectPolicyFunc,
  }

  response, err := client.Get(rest_url_for("/databases"))
  if err != nil {
    return nil, err
  }
  fmt.Println(response)

  request, err := http.NewRequest("GET", rest_url_for("/databases"), nil)
  request.Header.Add("Authorization", "Bearer ") // + oauth_token)

  return nil, err
}
