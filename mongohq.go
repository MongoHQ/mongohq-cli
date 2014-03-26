package main

import (
  //"fmt"
  "os"
  "github.com/codegangsta/cli"
  "github.com/MongoHQ/controllers"  // MongoHQ CLI functions
)

func main() {
  app := cli.NewApp()
  app.Name = "mongohq"
  app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
  app.Before = controllers.RequireAuth
  app.Commands = []cli.Command{
    {
      Name:      "databases",
      Usage:     "list databases",
      Action: func(c *cli.Context) {
        controllers.Databases()
      },
    },
    {
      Name:      "databases:info (pending)",
      Usage:     "information on database",
      Flags:     []cli.Flag {
        cli.StringFlag { "db", "database-name", ""},
      },
      Action: func(c *cli.Context) {
        println("Pending")
      },
    },
    {
      Name:      "deployments",
      Usage:     "list deployments",
      Action: func(c *cli.Context) {
        controllers.Deployments()
      },
    },
    {
      Name:      "deployments:info (pending)",
      Usage:     "information on deployment",
      Flags:     []cli.Flag {
        cli.StringFlag { "dp", "host:port", ""},
      },
      Action: func(c *cli.Context) {
        println("Pending")
      },
    },
    {
      Name:      "deployments:mongostat",
      Usage:     "realtime mongostat",
      Flags:     []cli.Flag {
        cli.StringFlag{"deployment,dp", "<bson_id>", "deployment id"},
        cli.StringFlag{"database,db", "<string>", "database name"},
      },
      Action: func(c *cli.Context) {
        if c.String("deployment") != "<bson_id>" {
          controllers.MongoStat(c.String("deployment"), c.String("database"))
        } else {
          println("Deployment is required")
          os.Exit(1)
        }
      },
    },
    {
      Name:      "deployments:logs (pending)",
      Usage:     "tail logs",
      Action: func(c *cli.Context) {
        println("Pending")
      },
    },
    {
      Name:      "logout",
      Usage:     "remove stored auth",
      Action:    func(c *cli.Context) {
        controllers.Logout()
      },
    },
  }

  app.Run(os.Args)
}
