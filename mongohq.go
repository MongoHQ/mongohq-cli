package main

import (
  //"fmt"
  "os"
  "github.com/codegangsta/cli"
  "github.com/MongoHQ/mongohq-cli"
  "github.com/MongoHQ/controllers"  // MongoHQ CLI functions
)

func main() {
  app := cli.NewApp()
  app.Name = "mongohq"
  app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
  app.Before = controllers.RequireAuth
  app.Version = mongohq_cli.Version()
  app.Commands = []cli.Command{
    {
      Name:      "databases",
      Usage:     "list databases",
      Action: func(c *cli.Context) {
        controllers.Databases()
      },
    },
    {
      Name:      "databases:info",
      Usage:     "information on database",
      Flags:     []cli.Flag {
        cli.StringFlag { "database,db", "<database name>", ""},
      },
      Action: func(c *cli.Context) {
        controllers.Database(c.String("database"))
      },
    },
    {
      Name:      "databases:users",
      Usage:     "users on a database",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dp", "<deployment id>", ""},
        cli.StringFlag { "database,db", "<database name>", ""},
      },
      Action: func(c *cli.Context) {
        controllers.DatabaseUsers(c.String("deployment"), c.String("database"))
      },
    },
    {
      Name:      "databases:users:create",
      Usage:     "add user to a database",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dp", "<deployment id>", ""},
        cli.StringFlag { "database,db", "<database name>", ""},
        cli.StringFlag { "username,u", "<user>", ""},
      },
      Action: func(c *cli.Context) {
        controllers.DatabaseCreateUser(c.String("deployment"), c.String("database"), c.String("username"))
      },
    },
    {
      Name:      "databases:users:remove",
      Usage:     "remove user from database",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dp", "<deployment id>", ""},
        cli.StringFlag { "database,db", "<database name>", ""},
        cli.StringFlag { "username,u", "<user>", ""},
      },
      Action: func(c *cli.Context) {
        controllers.DatabaseRemoveUser(c.String("deployment"), c.String("database"), c.String("username"))
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
      Name:      "deployments:info",
      Usage:     "information on deployment",
      Flags:     []cli.Flag {
        cli.StringFlag { "deployment,dp", "host:port", ""},
      },
      Action: func(c *cli.Context) {
        controllers.Deployment(c.String("deployment"))
      },
    },
    {
      Name:      "deployments:mongostat",
      Usage:     "realtime mongostat",
      Flags:     []cli.Flag {
        cli.StringFlag{"deployment,dp", "<bson_id>", "deployment id"},
      },
      Action: func(c *cli.Context) {
        if c.String("deployment") != "<bson_id>" {
          controllers.DeploymentMongoStat(c.String("deployment"))
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
      Name:      "deployments:oplog",
      Usage:     "tail oplog",
      Flags:     []cli.Flag {
        cli.StringFlag{"deployment,dp", "<bson_id>", "deployment id"},
      },
      Action: func(c *cli.Context) {
        if c.String("deployment") != "<bson_id>" {
          controllers.DeploymentOplog(c.String("deployment"))
        } else {
          println("Deployment is required")
          os.Exit(1)
        }

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
