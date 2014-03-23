package main

import (
  //"fmt"
  "os"
  "github.com/codegangsta/cli"
  "github.com/MongoHQ/mongohq_cli"  // MongoHQ CLI functions
)

func main() {
  app := cli.NewApp()
  app.Name = "mongohq"
  app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
  app.Before = mongohq_cli.RequireAuth
  app.Commands = []cli.Command{
    {
      Name:      "databases",
      Usage:     "list databases",
      Action: func(c *cli.Context) {
        mongohq_cli.Databases()
      },
    },
    {
      Name:      "databases:info",
      Usage:     "information on database",
      Flags:     []cli.Flag {
        cli.StringFlag { "db", "database-name", ""},
      },
      Action: func(c *cli.Context) {
        println("added task: ", c.Args().First())
      },
    },
    {
      Name:      "deployments",
      Usage:     "list deployments",
      Action: func(c *cli.Context) {
        println("added task: ", c.Args().First())
      },
    },
    {
      Name:      "deployments:info",
      Usage:     "information on deployment",
      Flags:     []cli.Flag {
        cli.StringFlag { "dp", "host:port", ""},
      },
      Action: func(c *cli.Context) {
        println("added task: ", c.Args().First())
      },
    },
    {
      Name:      "deployments:mongostat",
      Usage:     "realtime mongostat",
      Action: func(c *cli.Context) {
        println("added task: ", c.Args().First())
      },
    },
    {
      Name:      "deployments:logs",
      Usage:     "tail logs",
      Action: func(c *cli.Context) {
        println("added task: ", c.Args().First())
      },
    },
    {
      Name:      "logout",
      Usage:     "remove stored auth",
      Action:    func(c *cli.Context) {
        mongohq_cli.Logout()
      },
    },
  }

  app.Run(os.Args)
}
