package controllers

import (
  "fmt"
  "github.com/MongoHQ/api"
  "os"
)

func Regions() {
  regions, err := api.GetRegions(OauthToken)

  if err != nil {
    fmt.Println("Error returning regions: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("== Regions")
  for _, region := range regions {
    fmt.Println(region)
  }
}


