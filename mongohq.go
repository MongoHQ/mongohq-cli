package main

import (
	"fmt"
	"controllers" // MongoHQ CLI functions
	"mongohq-cli"
	"github.com/codegangsta/cli"
	"os"
)

func requireArguments(c *cli.Context, argumentsSlice []string, errorMessages []string) {
	err := false

	for _, argument := range argumentsSlice {
		if !c.IsSet(argument) {
			err = true
			fmt.Println("--" + argument + " is required")
		}
	}

	if err {
		fmt.Println("\nMissing arguments, for more information, run: mongohq " + c.Command.Name + " --help\n")
		for _, errorMessage := range errorMessages {
			fmt.Println(errorMessage)
		}
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "mongohq"
	app.Usage = "Allow MongoHQ interaction from the commandline (enables awesomeness)"
	app.Before = controllers.RequireAuth
	app.Version = mongohq_cli.Version()
	app.Commands = []cli.Command{
		{
			Name:  "backups",
			Usage: "list backups with optional filters",
			Flags: []cli.Flag{
				cli.StringFlag{"database,db", "<string>", "(optional) database to list backups for"},
				cli.StringFlag{"deployment,dep", "<string>", "(optional) deployment to list backups for"},
			},
			Action: func(c *cli.Context) {
				filter := map[string]string{}
				if c.IsSet("database") {
					filter["database"] = c.String("database")
				}
				if c.IsSet("deployment") {
					filter["deployment"] = c.String("deployment")
				}
				controllers.Backups(filter)
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
				controllers.Backup(c.String("backup"))
			},
		},
		{
			Name:  "backups:restore",
			Usage: "restore backup to a new database",
			Flags: []cli.Flag{
				cli.StringFlag{"backup,b", "<string>", "file name of backup"},
				cli.StringFlag{"source-database,source", "<string>", "original database name"},
				cli.StringFlag{"destination-database,destination", "<string>", "new database name"},
			},
			Action: func(c *cli.Context) {
				requireArguments(c, []string{"backup", "source-database", "destination-database"}, []string{})
				controllers.RestoreBackup(c.String("backup"), c.String("source-database"), c.String("destination-database"))
			},
		},
		{
			Name:  "databases",
			Usage: "list databases",
			Action: func(c *cli.Context) {
        controllers.Databases()
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
				controllers.CreateDatabase(c.String("deployment"), c.String("database"))
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
				controllers.Database(c.String("database"))
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
				controllers.RemoveDatabase(c.String("database"), c.Bool("force"))
			},
		},
		{
			Name:  "deployments",
			Usage: "list deployments",
			Action: func(c *cli.Context) {
				controllers.Deployments()
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
				controllers.CreateDeployment(c.String("deployment"), c.String("database"), c.String("region"))
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
				controllers.Deployment(c.String("deployment"))
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
				controllers.DeploymentRename(c.String("deployment"), c.String("name"))
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
				controllers.DeploymentMongoStat(c.String("deployment"))
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
				controllers.HistoricalLogs(c.String("deployment"))
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
				//controllers.DeploymentOplog(c.String("deployment"))
			//},
		//},
		{
			Name:  "regions",
			Usage: "list available regions",
			Action: func(c *cli.Context) {
				controllers.Regions()
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
				controllers.DatabaseUsers(c.String("deployment"), c.String("database"))
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
				controllers.DatabaseCreateUser(c.String("deployment"), c.String("database"), c.String("username"))
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
				controllers.DatabaseRemoveUser(c.String("deployment"), c.String("database"), c.String("username"))
			},
		},
		{
			Name:  "logout",
			Usage: "remove stored auth",
			Action: func(c *cli.Context) {
				controllers.Logout()
			},
		},
	}

	app.Run(os.Args)
}
