package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
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
		loginController.Api.Config = getConfig()

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
				cli.StringFlag{"account,a", "<string>", "optional account slug; if included, will run accounts:info"},
			},
			Description: `
   List the slugs for all accounts which you have permission
   To change the default account, see the "config:account" command.
      `,
			Action: func(c *cli.Context) {
				if c.String("account") == "<string>" {
					controller.ListAccounts()
				} else {
					controller.ShowAccount(c.String("account"))
				}
			},
		},
		{
			Name:  "accounts:info",
			Usage: "account information",
			Flags: []cli.Flag{
				cli.StringFlag{"account,a", "<string>", "account slug"},
			},
			Description: `
   More detail about a particular account, including name, slug, owner, and users.
      `,
			Action: func(c *cli.Context) {
				controller.ShowAccount(c.String("account"))
			},
		},
		{
			Name:  "backups",
			Usage: "list backups with optional filters",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "optional deployment filter for backups"},
				cli.StringFlag{"backup,b", "<string>", "optional backup name; if included, will run backups:info"},
			},
			Description: `
   Lists the backups associated with your account or deployment.

   To see a list of all backups on your account, including those from deleted
   deployments, omit the deployment argument.

   To see a list of all backups on a single deployment, include the name or
   id of the intended deployment.
      `,
			Action: func(c *cli.Context) {
				if c.String("backup") == "<string>" {
					controller.ListBackups(c.String("deployment"))
				} else {
					controller.ShowBackup(c.String("backup"))
				}
			},
		},
		{
			Name:  "backups:create",
			Usage: "create an on-demand backup for a deployment",
			Description: `
   Queues an on-demand backup for a deployment.  To read more about this feature,
   see http://docs.mongohq.com/backups/elastic-deployments.html#on-demand-backups.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment name"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment"}, []string{})
				controller.CreateBackup(c.String("deployment"))
			},
		},
		{
			Name:  "backups:info",
			Usage: "information on backup",
			Flags: []cli.Flag{
				cli.StringFlag{"backup,b", "<string>", "file name of backup"},
			},
			Description: `
   More detail about a particular backup, including deployment, databases, creation time, type, size, and download link.
      `,
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
   deployment will be created in the same datacenter with the same version
   as the source database.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "new deployment name"},
				cli.StringFlag{"backup,b", "<string>", "file name of backup"},
				cli.StringFlag{"source-database,source", "<string>", "original database name"},
				cli.StringFlag{"destination-database,destination", "<string>", "new database name"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"deployment", "backup", "source-database", "destination-database"}, []string{})
				controller.RestoreBackup(c.String("backup"), c.String("deployment"), c.String("source-database"), c.String("destination-database"))
			},
		},
		{
			Name:  "config:account",
			Usage: "set a default account context",
			Flags: []cli.Flag{
				cli.StringFlag{"account,a", "<string>", "slug for default account"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"account"}, []string{})
				controller.SetConfigAccount(c.String("account"))
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
				cli.StringFlag{"deployment,dep", "<string>", " deployment containing database"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"database", "deployment"}, []string{})
				controller.ShowDatabase(c.String("deployment"), c.String("database"))
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
			Description: `
   List the slugs for all deployments.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "optional deployment name; if included runs deployments:info"},
			},
			Action: func(c *cli.Context) {
				if c.String("deployment") == "<string>" {
					controller.ListDeployments()
				} else {
					controller.ShowDeployment(c.String("deployment"))
				}
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
			Name:  "whoami",
			Usage: "display effective user",
			Action: func(c *cli.Context) {
				controller.CurrentUser()
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
			Name:  "update",
			Usage: "script to update the MongoHQ CLI binary",
			Action: func(c *cli.Context) {
				fmt.Println("To update, run: `curl https://mongohq-cli.s3.amazonaws.com/install.sh | sh`")
			},
		},
	}

	app.Run(os.Args)
}
