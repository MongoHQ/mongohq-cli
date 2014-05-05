package main

import (
  "fmt"
  "os"
)

func (c *Controller) ListRegions() {
  providersRegions, err := c.Api.GetRegions()

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


