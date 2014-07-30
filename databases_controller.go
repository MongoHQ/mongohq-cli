package main

import (
	"fmt"
)

func (c *Controller) ListDatabases() {
	databases, err := c.Api.GetDatabases()

	if err != nil {
		fmt.Println("Error retrieving databases: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== My Databases")
	for _, database := range databases {
		fmt.Println(database.Name)
	}
}

func (c *Controller) ShowDatabase(deploymentName, databaseName string) {
	database, err := c.Api.GetDatabase(deploymentName, databaseName)

	if err != nil {
		fmt.Println("Error retrieiving database: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== " + database.Name)
	fmt.Println(" name       : " + database.Name)
	fmt.Println(" plan       : " + database.Plan)
	fmt.Println(" status     : " + database.Status)
	fmt.Println(" deployment : " + deploymentName)

	if database.Status == "running" {
		users, err := c.Api.GetDatabaseUsers(deploymentName, databaseName)

		if err != nil {
			fmt.Println(" == Error returning database users: " + err.Error())
		} else {
			fmt.Println(" == users")
			for _, user := range users {
				fmt.Println("  " + user.Username)
			}
		}

		stats, err := c.Api.GetDatabaseStats(database)

		if err != nil {
			fmt.Println(" == Error returning database stats: " + err.Error())
		} else {
			fmt.Println(" == database usage stats per host")
			for host, stat := range stats {
				fmt.Println("  == " + host)
				fmt.Println("   dataSize:  " + prettySize(float64(stat.DataSize)))
				fmt.Println("   indexSize: " + prettySize(float64(stat.IndexSize)))
				fmt.Println("   fileSize:  " + prettySize(float64(stat.FileSize)))
			}
		}
	}
}

func (c *Controller) DeleteDatabase(deploymentSlug, databaseName string, force bool) {
	if !force {
		confirmDatabaseName := prompt("To confirm, type the name of the database to be deleted")

		if databaseName != confirmDatabaseName {
			fmt.Println("Confirmation of database name is incorrect.")
			cliOSExit()
			return
		}
	}

	err := c.Api.RemoveDatabase(deploymentSlug, databaseName)

	if err != nil {
		fmt.Println("Error removing database: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("Removed database named: " + databaseName)
}

func (c *Controller) CreateDatabase(deploymentName, databaseName string) {
	database, err := c.Api.CreateDatabase(deploymentName, databaseName)

	if err != nil {
		fmt.Println("Error creating database: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== " + database.Name)
	fmt.Println(" status     :" + database.Status)
	fmt.Println(" deployment :" + deploymentName)
}

func (c *Controller) ListDatabaseUsers(deploymentId, databaseName string) {
	databaseUsersSlice, err := c.Api.GetDatabaseUsers(deploymentId, databaseName)

	if err != nil {
		fmt.Println("Error retrieiving database users: " + err.Error())
		cliOSExit()
		return
	} else {
		fmt.Println("== Users for database " + databaseName)
		for _, databaseUser := range databaseUsersSlice {
			fmt.Println("  " + databaseUser.Username)
		}
	}
}

func (c *Controller) CreateDatabaseUser(deploymentId, databaseName, username, suppliedPassword string) {
	var password string
	var err error

	if suppliedPassword == "<string>" {
		password, err = safeGetPass("Password (typing will be hidden): ")

		if err != nil {
			fmt.Println("Error returning password.  We may not be compliant with your system yet.  Please send us a message telling us about your system to support@mongohq.com.")
			cliOSExit()
			return
		}

		confirmedPassword, _ := safeGetPass("Confirm password: ")

		if password != confirmedPassword {
			fmt.Println("Password confirmation failed.")
			cliOSExit()
			return
		}
	} else {
		password = suppliedPassword
	}

	_, err = c.Api.CreateDatabaseUser(deploymentId, databaseName, username, password)

	if err != nil {
		fmt.Println("Error creating database user: " + err.Error())
		cliOSExit()
		return
	}
	fmt.Println("User " + username + " created.")
}

func (c *Controller) DeleteDatabaseUser(deploymentId, databaseName, username string) {
	_, err := c.Api.RemoveDatabaseUser(deploymentId, databaseName, username)

	if err != nil {
		fmt.Println("Error removing database user: " + err.Error())
		cliOSExit()
		return
	}
	fmt.Println("User " + username + " removed.")
}
