package controllers

import (
  "fmt"
)

func prompt(text string) (string) {
  var response string
  var err error

  print(text + ": ")
  _, err = fmt.Scanln(&response)
  if err != nil {
    fmt.Println("Error: ", err)
  }
  return response
}
