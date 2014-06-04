package main

import (
	"github.com/codegangsta/cli"
	"os"
  "fmt"
)

var api *Api
var controller Controller

func main() {
  loginController := new(LoginController)

	app := cli.NewApp()
	app.Name = "mongohq"
	app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
	app.Before = func(c *cli.Context) error {
    loginController.RequireAuth(c) // Exits process if auth fails
    loginController.Api.Defaults = getDefaults()
    controller = Controller{Api: loginController.Api}
    return nil
  }
  app.CommandNotFound = findClosestCommand
	app.Version = Version()
	app.Commands = []cli.Command{
		{
			Name:  "accounts",
			Usage: "list accounts",
			Flags: []cli.Flag{
			},
			Action: func(c *cli.Context) {
				controller.ListAccounts()
			},
		},
		{
      Name:  "accounts:info",
			Usage: "account information",
			Flags: []cli.Flag{
				cli.StringFlag{"account,a", "<string>", "account slug"},
			},
			Action: func(c *cli.Context) {
				controller.ShowAccount(c.String("account"))
			},
		},
		{
			Name:  "backups",
			Usage: "list backups with optional filters",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "(optional) deployment to list backups for"},
			},
			Action: func(c *cli.Context) {
				controller.ListBackups(c.String("deployment"))
			},
		},
		{
			Name:  "backups:info",
			Usage: "information on backup",
			Flags: []cli.Flag{
				cli.StringFlag{"backup,b", "<string>", "file name of backup"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"backup"}, []string{})
				controller.ShowBackup(c.String("backup"))
			},
		},
		{
			Name:  "backups:restore",
			Usage: "restore backup to a new database",
      Description: `
      Restores a backup of a database to a new, fresh deployment. The new
      deployment will be created in the same datacenter as the source database.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "new deployment name"},
				cli.StringFlag{"backup,b", "<string>", "file name of backup"},
				cli.StringFlag{"source-database,source", "<string>", "original database name"},
				cli.StringFlag{"destination-database,destination", "<string>", "new database name"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "backup", "source-database", "destination-database"}, []string{})
				controller.RestoreBackup(c.String("backup"), c.String("source-database"), c.String("destination-database"))
			},
		},
    {
      Name: "config:account",
      Usage: "set a default account context",
      Flags: []cli.Flag{
        cli.StringFlag{"account,a", "<string>", "slug for default account"},
      },
      Action: func(c *cli.Context) {
				requireArguments(c, []string{"account"}, []string{})
        controller.SetDefaultAccount(c.String("account"))
      },
    },
		{
			Name:  "databases:create",
			Usage: "create database on an existing deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment to create database on"},
				cli.StringFlag{"database,db", "<string>", "new database to create"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "database"}, []string{})
				controller.CreateDatabase(c.String("deployment"), c.String("database"))
			},
		},
		{
			Name:  "databases:info",
			Usage: "information on database",
			Flags: []cli.Flag{
				cli.StringFlag{"database,db", "<string>", " database for more information"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"database"}, []string{})
				controller.ShowDatabase(c.String("database"))
			},
		},
		{
			Name:  "databases:remove",
			Usage: "remove database",
			Flags: []cli.Flag{
				cli.StringFlag{"database,db", "<string>", "database to remove"},
				cli.BoolFlag{"force,f", "delete without confirmation"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"database"}, []string{})
				controller.DeleteDatabase(c.String("database"), c.Bool("force"))
			},
		},
		{
			Name:  "deployments",
			Usage: "list deployments",
			Action: func(c *cli.Context) {
				controller.ListDeployments()
			},
		},
		{
			Name:  "deployments:create",
			Usage: "create a new Elastic Deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"database,db", "<string>", "new database name"},
				cli.StringFlag{"deployment,dep", "<string>", "new deployment name"},
				cli.StringFlag{"region,r", "<string>", "region of deployment (for list of regions, run 'mongohq regions')"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "database", "region"}, []string{})
				controller.CreateDeployment(c.String("deployment"), c.String("database"), c.String("region"))
			},
		},
		{
			Name:  "deployments:info",
			Usage: "information on deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for more information"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment"}, []string{})
				controller.ShowDeployment(c.String("deployment"))
			},
		},
		{
			Name:  "deployments:rename",
			Usage: "rename a deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for more information"},
				cli.StringFlag{"name,n", "<string>", "new name for deployment"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "name"}, []string{})
				controller.RenameDeployment(c.String("deployment"), c.String("name"))
			},
		},
		{
			Name:  "mongostat",
			Usage: "realtime mongostat",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for watching mongostats"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment"}, []string{})
				controller.DeploymentMongoStat(c.String("deployment"))
			},
		},
		{
			Name:  "logs",
			Usage: "query historical logs",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for querying logs"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment"}, []string{})
				controller.HistoricalLogs(c.String("deployment"))
			},
		},
		//{
			//Name:  "deployments:oplog",
			//Usage: "tail oplog",
			//Flags: []cli.Flag{
				//cli.StringFlag{"deployment,dep", "<string>", "deployment to tail the oplog"},
			//},
			//Action: func(c *cli.Context) {
				//requireArguments("deployments:oplog", c, []string{"deployment"}, []string{})
				//controller.DeploymentOplog(c.String("deployment"))
			//},
		//},
		{
			Name:  "regions",
			Usage: "list available regions",
			Action: func(c *cli.Context) {
				controller.ListRegions()
			},
		},
		{
			Name:  "users",
			Usage: "list users on a database",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment id the database is on"},
				cli.StringFlag{"database,db", "<string>", "database to list users"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "database"}, []string{})
				controller.ListDatabaseUsers(c.String("deployment"), c.String("database"))
			},
		},
		{
			Name:  "users:create",
			Usage: "add user to a database",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment id the database is on"},
				cli.StringFlag{"database,db", "<string>", "atabase name to create the user on"},
				cli.StringFlag{"username,u", "<string>", "user to create"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "database", "username"}, []string{})
				controller.CreateDatabaseUser(c.String("deployment"), c.String("database"), c.String("username"))
			},
		},
		{
			Name:  "users:remove",
			Usage: "remove user from database",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment id the database is on"},
				cli.StringFlag{"database,db", "<string>", "database name to remove the user from"},
				cli.StringFlag{"username,u", "<string>", "user to remove from the deployment"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "database", "username"}, []string{})
				controller.DeleteDatabaseUser(c.String("deployment"), c.String("database"), c.String("username"))
			},
		},
		{
			Name:  "logout",
			Usage: "remove stored auth",
			Action: func(c *cli.Context) {
				loginController.Logout()
			},
		},
    {
      Name: "update",
      Usage: "script to update the MongoHQ CLI binary",
      Action: func(c *cli.Context) {
        fmt.Println("To update, run: `curl https://mongohq-cli.s3.amazonaws.com/install.sh | sh`")
      },
    },
	}

	app.Run(os.Args)
}
