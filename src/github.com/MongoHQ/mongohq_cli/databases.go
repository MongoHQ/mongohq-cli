package mongohq_cli

import (
  "fmt"
  "github.com/MongoHQ/mongohq_api"
)

func Databases() {
  databases := mongohq_api.Databases
  fmt.Println(databases)
}
