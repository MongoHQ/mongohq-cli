package main

import (
  "fmt"
  "os"
)

func (c *Controller) ListAccounts() {
  accountsSlice, err := c.Api.GetAccounts()

  if err != nil {
    fmt.Println("Error retreiving accounts: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("== Accounts")
  for _, account := range accountsSlice {
    fmt.Println(account.Slug)
  }
}

func (c *Controller) ShowAccount(slug string) {
  account, err := c.Api.GetAccount(slug)

  if err != nil {
    fmt.Println("Error retreiving accounts: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("== " + slug)
  fmt.Println(" slug:    " + account.Slug)
  fmt.Println(" name:    " + account.Name)
  fmt.Println(" owner:   " + account.OwnerName)
}
