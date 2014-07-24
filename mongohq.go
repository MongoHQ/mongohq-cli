package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/peterh/liner"
	"os"
	"strings"
)

var api *Api
var controller Controller
var historyfn="/tmp/.mongohq_history"

func main() {
	loginController := new(LoginController)

	app := cli.NewApp()
	app.Name = "mongohq"
	app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
	app.Before = func(c *cli.Context) error {
		loginController.Api = &Api{UserAgent: "MongoHQ-CLI " + Version()}
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
List the slugs for all accounts which you have permission To change the default account, see the "config:account" command.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()

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
More detail about a particular account, including name, slug, owner, and account users.

These account users are different than database users, and cannot be used to directly access a database.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireArguments(c, []string{"account"}, []string{})
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

To see a list of all backups on your account, including those from deleted deployments, omit the deployment argument.

To see a list of all backups on a single deployment, include the name or id of the intended deployment using the deployment argument.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				if c.String("backup") == "<string>" {
					if c.String("deployment") == "<string>" {
						controller.ListBackups()
					} else {
						controller.ListBackupsForDeployment(c.String("deployment"))
					}
				} else {
					controller.ShowBackup(c.String("backup"))
				}
			},
		},
		{
			Name:  "backups:create",
			Usage: "create an on-demand backup for a deployment",
			Description: `
Queues an on-demand backup for a deployment.  To read more about this feature, see http://docs.mongohq.com/backups/elastic-deployments.html#on-demand-backups.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment name"},
			},
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

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
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"backup"}, []string{})
				controller.ShowBackup(c.String("backup"))
			},
		},
		{
			Name:  "backups:restore",
			Usage: "restore backup to a new database",
			Description: `
Restores a backup of a database to a new, fresh deployment. The new deployment will be created in the same datacenter with the same version as the source database.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "new deployment name"},
				cli.StringFlag{"backup,b", "<string>", "file name of backup"},
				cli.StringFlag{"source-database,source", "<string>", "original database name"},
				cli.StringFlag{"destination-database,destination", "<string>", "new database name"},
			},
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

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
			Description: `
Set a default account so the account flag is not required for each command.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()

				requireArguments(c, []string{"account"}, []string{})
				controller.SetConfigAccount(c.String("account"))
			},
		},
		{
			Name:      "databases:create",
			ShortName: "db:create",
			Usage:     "create database on an existing deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment to create database on"},
				cli.StringFlag{"database,db", "<string>", "new database to create"},
			},
			Description: `
Create a new database on an existing deployment.  If you are looking to create a new database on a new deployment, see the deployments:create command.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment", "database"}, []string{})
				controller.CreateDatabase(c.String("deployment"), c.String("database"))
			},
		},
		{
			Name:      "databases:info",
			ShortName: "db:info",
			Usage:     "information on database",
			Flags: []cli.Flag{
				cli.StringFlag{"database,db", "<string>", " database for more information"},
				cli.StringFlag{"deployment,dep", "<string>", " deployment containing database"},
			},
			Description: `
More detail on a particular database, including name, status, and stats.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"database", "deployment"}, []string{})
				controller.ShowDatabase(c.String("deployment"), c.String("database"))
			},
		},
		{
			Name:      "databases:remove",
			ShortName: "db:remove",
			Usage:     "remove database",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment"},
				cli.StringFlag{"database,db", "<string>", "database to remove"},
				cli.BoolFlag{"force,f", "delete without confirmation"},
			},
			Description: `
Deletes a database from a deployment.  If this is the last database on the deployment, the deployment will also be deleted.

You will be asked to verify the database name on delete, unless including the force argument.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"database", "deployment"}, []string{})
				controller.DeleteDatabase(c.String("deployment"), c.String("database"), c.Bool("force"))
			},
		},
		{
			Name:      "deployments",
			ShortName: "dep",
			Usage:     "list deployments",
			Description: `
List the slugs for all deployments.
      `,
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "optional deployment name; if included runs deployments:info"},
			},
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				if c.String("deployment") == "<string>" {
					controller.ListDeployments()
				} else {
					controller.ShowDeployment(c.String("deployment"))
				}
			},
		},
		{
			Name:      "deployments:create",
			ShortName: "dep:create",
			Usage:     "create a new Elastic Deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"database,db", "<string>", "new database name"},
				cli.StringFlag{"deployment,dep", "<string>", "new deployment name"},
				cli.StringFlag{"location,l", "<string>", "location of deployment (for list of locations, run 'mongohq locations')"},
			},
			Description: `
Creates an elastic deployment on the MongoHQ platform. Stick with me here: it will create a new database on a new deployment at location you specify.  The deployment is a Replica Set and the database is the logical MongoDB database. You can find a list of locations by running "mongohq locations".
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment", "database", "location"}, []string{})
				controller.CreateDeployment(c.String("deployment"), c.String("database"), c.String("location"))
			},
		},
		{
			Name:      "deployments:info",
			ShortName: "dep:info",
			Usage:     "information on deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for more information"},
			},
			Description: `
More detail about a particular deployment, including plan, status, location, current primary, members, version, and a list of databases.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment"}, []string{})
				controller.ShowDeployment(c.String("deployment"))
			},
		},
		{
			Name:      "deployments:rename",
			ShortName: "dep:rename",
			Usage:     "rename a deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for more information"},
				cli.StringFlag{"name,n", "<string>", "new name for deployment"},
			},
			Description: `
Sometime, you want a little more description about a deployment than an hex id.  Use this to create a deployment name (only allows alphanumeric characters and hyphens).

Immediately after making this change, you will need to reference the deployment by the new name.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment", "name"}, []string{})
				controller.RenameDeployment(c.String("deployment"), c.String("name"))
			},
		},
		{
			Name:      "deployments:remove",
			ShortName: "dep:remove",
			Usage:     "remove a deployment",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for more information"},
				cli.BoolFlag{"force,f", "delete without confirmation"},
			},
			Description: `
Deletes a deployment.  Requires confirmation because this is a very destructive action, particularly for data.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment"}, []string{})
				controller.DeleteDeployment(c.String("deployment"), c.Bool("force"))
			},
		},
		{
			Name:  "logs",
			Usage: "query historical logs",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for querying logs"},
				cli.StringFlag{"regexp,e", "<string>", "regexp for log searches"},
				cli.StringFlag{"search,s", "<string>", "exact search term for log searches"},
				cli.StringFlag{"exclude,v", "<string>", "exclude search term for log searches"},
			},
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment"}, []string{})
				controller.HistoricalLogs(c.String("deployment"), c.String("search"), c.String("exclude"), c.String("regexp"))
			},
		},
		{
			Name:  "locations",
			Usage: "list available locations",
			Description: `
List the current locations available for MongoHQ deployments.  Used with both new deployments and restoring databases from backups.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				controller.ListLocations()
			},
		},
		{
			Name:  "mongostat",
			Usage: "realtime mongostat",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment for watching mongostats"},
			},
			Description: `
A streaming output of usage statistics for your database.  This is a very good first step when you are looking for performance characteristics on your database.  The usage stats include:

 * Operational stats: inserts, queries, updates, deletes, getmores, commands per second
 * Memory usage: physical and virtual usage, with page swaps (i.e. faults) / second
 * Database behavior: flushes, locked percentage, queued reads and writes
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment"}, []string{})
				controller.DeploymentMongoStat(c.String("deployment"))
			},
		},
		{
			Name:  "users",
			Usage: "list users on a database",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment id the database is on"},
				cli.StringFlag{"database,db", "<string>", "database to list users"},
			},
			Description: `
List a databases' users.  These users are used to authenticate against a database.

These are different than account users, which are used to authentication against the MongoHQ service.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment", "database"}, []string{})
				controller.ListDatabaseUsers(c.String("deployment"), c.String("database"))
			},
		},
		{
			Name:  "users:create",
			Usage: "add user to a database",
			Flags: []cli.Flag{
				cli.StringFlag{"deployment,dep", "<string>", "deployment id the database is on"},
				cli.StringFlag{"database,db", "<string>", "database name to create the user on"},
				cli.StringFlag{"username,u", "<string>", "user to create"},
				cli.StringFlag{"password,p", "<string>", "optional password for user; will prompt if omitted"},
			},
			Description: `
Add a new user to a database. With this user, you will be able to authenticate against the database. If a password is not provided, it will be prompted.

If the user already exists, this command will update the password for the user.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment", "database", "username"}, []string{})
				controller.CreateDatabaseUser(c.String("deployment"), c.String("database"), c.String("username"), c.String("password"))
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
			Description: `
Removes a database user from a database.  If your applications are connecting with this user, they will not be able to create new connections.

This user action is against database users used for authentication against a database.  It is different than account users.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()
				requireAccount(loginController.Api)

				requireArguments(c, []string{"deployment", "database", "username"}, []string{})
				controller.DeleteDatabaseUser(c.String("deployment"), c.String("database"), c.String("username"))
			},
		},
		{
			Name:  "whoami",
			Usage: "display effective user",
			Description: `
Just a simple command to tell you which account user you are currently acting as.
      `,
			Action: func(c *cli.Context) {
				loginController.RequireAuth()

				controller.CurrentUser()
			},
		},
		{
			Name:  "logout",
			Usage: "remove stored auth",
			Description: `
Removes authentication information from the MongoHQ CLI on this machine, and sends a kill command to the oauth token used for authentication.
      `,
			Action: func(c *cli.Context) {
				loginController.Logout()
			},
		},
		{
			Name:  "update",
			Usage: "script to update the MongoHQ CLI binary",
			Description: `
To update, run:

  curl https://mongohq-cli.s3.amazonaws.com/install.sh | sh
      `,
			Action: func(c *cli.Context) {
				fmt.Println("To update, run: `curl https://mongohq-cli.s3.amazonaws.com/install.sh | sh`")
			},
		},
	}

	if len(os.Args)==1 {
		term:=liner.NewLiner()
		var prompt="> "
		if f,err := os.Open(historyfn); err==nil {
			term.ReadHistory(f)
			f.Close()
		}
		for {
			line,err:=term.Prompt(prompt)
			if err!=nil {
				break
			}
			if line=="exit" {
				break
			}
			term.AppendHistory(line)
			plines:=strings.Fields(line)
			plines=append(plines,"")
			copy(plines[1:],plines[0:])
			plines[0]="mongohq-repl"
			app.Run(plines)
		}
		if f,err :=os.Create(historyfn); err!=nil {
			fmt.Println("Error writing history file",err)
		} else {
			term.WriteHistory(f)
			f.Close()
		}
		os.Exit(0)
	} else {
		app.Run(os.Args)
	}
}
