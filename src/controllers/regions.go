package controllers

import (
  "fmt"
  "api"
  "os"
)

func Regions() {
  providersRegions, err := api.GetRegions(OauthToken)

  if err != nil {
    fmt.Println("Error returning regions: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("== providers & regions")
  for provider, regions := range providersRegions {
    fmt.Println("  " + provider)
    for _, region := range regions {
      fmt.Println("   - " + region)
    }
  }
}


