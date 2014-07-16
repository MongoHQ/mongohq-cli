package main

import (
	"fmt"
	"os"
)

func (c *Controller) ListDatabases() {
	databases, err := c.Api.GetDatabases()

	if err != nil {
		fmt.Println("Error retrieving databases: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("=== My Databases")
	for _, database := range databases {
		fmt.Println(database.Name)
	}
}

func (c *Controller) ShowDatabase(deploymentName, databaseName string) {
	database, err := c.Api.GetDatabase(deploymentName, databaseName)

	if err != nil {
		fmt.Println("Error retrieiving database: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("=== " + database.Name)
	fmt.Println(" id:         " + database.Id)
	fmt.Println(" name:       " + database.Name)
	fmt.Println(" plan:       " + database.Plan)
	fmt.Println(" status:     " + database.Status)
	fmt.Println(" deployment: " + deploymentName)
}

func (c *Controller) DeleteDatabase(databaseName string, force bool) {
	if !force {
		confirmDatabaseName := prompt("To confirm, type the name of the database to be deleted")

		if databaseName != confirmDatabaseName {
			fmt.Println("Confirmation of database name is incorrect.")
			os.Exit(1)
		}
	}

	err := c.Api.RemoveDatabase(databaseName)

	if err != nil {
		fmt.Println("Error removing database: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Removed database named: " + databaseName)
}

func (c *Controller) CreateDatabase(deploymentName, databaseName string) {
	database, err := c.Api.CreateDatabase(deploymentName, databaseName)

	if err != nil {
		fmt.Println("Error creating database: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("=== " + database.Name)
	fmt.Println(" status:        " + database.Status)
	fmt.Println(" deployment:    " + deploymentName)
}

func (c *Controller) ListDatabaseUsers(deploymentId, databaseName string) {
	databaseUsersSlice, err := c.Api.GetDatabaseUsers(deploymentId, databaseName)

	if err != nil {
		fmt.Println("Error retrieiving database users: " + err.Error())
		os.Exit(1)
	} else {
		fmt.Println("== Users for database " + databaseName)
		for _, databaseUser := range databaseUsersSlice {
			fmt.Println("  " + databaseUser.Username)
		}
	}
}

func (c *Controller) CreateDatabaseUser(deploymentId, databaseName, username string) {
	password := prompt("New user password")

	_, err := c.Api.CreateDatabaseUser(deploymentId, databaseName, username, password)

	if err != nil {
		fmt.Println("Error creating database user: " + err.Error())
		os.Exit(1)
	}
	fmt.Println("User " + username + " created.")
}

func (c *Controller) DeleteDatabaseUser(deploymentId, databaseName, username string) {
	_, err := c.Api.RemoveDatabaseUser(deploymentId, databaseName, username)

	if err != nil {
		fmt.Println("Error removing database users: " + err.Error())
		os.Exit(1)
	}
	fmt.Println("User " + username + " removed.")
}
