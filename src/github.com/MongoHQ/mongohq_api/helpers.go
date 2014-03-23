package mongohq_api

import (
  "net/http"
  "io/ioutil"
)

var oauth_client_id = "6fb9368538ef061ed73be71cc291e65b"
var oauth_secret    = "028d31d8ca253cc3004b3ae4470c21bb23c3011e2fc8b442ad72f259be7879ce5c66bfda4ff26d5a0ba8d23369ef3355ef4579f6e7a977ba933dc1a37fd2880c"

func rest_url_for(path string) (string) {
   return "https://dblayer-api.herokuapp.com" + path;
}

func rest_get(path string, oauthToken string) ([]byte, error) {
  client := &http.Client{}
  request, err := http.NewRequest("GET", rest_url_for(path), nil)
  request.Header.Add("Authorization", "Bearer " + oauthToken)
  response, err := client.Do(request)

  if err != nil {
    return nil, err
  } else {
    responseBody, _ := ioutil.ReadAll(response.Body)
    response.Body.Close()
    return responseBody, err
  }
}
