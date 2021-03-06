package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func (c *Controller) ListDeployments() {
	deployments, err := c.Api.GetDeployments()

	if err != nil {
		fmt.Println("Error retrieving deployments: " + err.Error())
	} else {
		fmt.Println("== My Deployments")
		for _, deployment := range deployments {
			fmt.Println(deployment.NameOrId())
		}
	}
}

func (c *Controller) ShowDeployment(deploymentId string) {
	deployment, err := c.Api.GetDeployment(deploymentId)

	if err != nil {
		fmt.Println("Error retrieving deployment: " + err.Error())
	} else {
		fmt.Println("== " + deployment.NameOrId())
		fmt.Println(" name            : " + deployment.NameOrId())
		fmt.Println(" plan            : " + deployment.Plan)
		fmt.Println(" status          : " + deployment.Status)
		fmt.Println(" location        : " + deployment.Location)
		fmt.Println(" current primary : " + deployment.CurrentPrimary)
		fmt.Println(" members         : " + strings.Join(deployment.Members, ","))
		fmt.Println(" version         : " + deployment.Version)

		if deployment.AllowMultipleDatabases {
			fmt.Println("  multiple databases?: true")
		}

		fmt.Println("  == Databases")
		for _, database := range deployment.Databases {
			fmt.Println("    " + database.Name)
		}
	}
}

func (c *Controller) RenameDeployment(deploymentId, name string) {
	_, err := c.Api.RenameDeployment(deploymentId, name)

	if err != nil {
		fmt.Println("Error renaming deployment: " + err.Error())
	} else {
		fmt.Println("Renamed deployment to " + name + ".  You will need to reference it by the new name.")
	}
}

func (c *Controller) CreateDeployment(deploymentName, databaseName, location string) {
	deployment, err := c.Api.CreateDeployment(deploymentName, databaseName, location)

	if err != nil {
		fmt.Println("Error creating deployment: " + err.Error())
	} else {
		fmt.Println("== Building deployment " + deploymentName + " with database " + databaseName + " in location " + location)

		c.pollNewDeployment(deployment)
	}
}

func (c *Controller) DeleteDeployment(deploymentName string, force bool) {
	if !force {
		confirmDeploymentName := prompt("To confirm, type the name of the deployment to be deleted")

		if deploymentName != confirmDeploymentName {
			fmt.Println("Confirmation of deployment name is incorrect.")
			cliOSExit()
			return
		}
	}

	err := c.Api.RemoveDeployment(deploymentName)

	if err != nil {
		fmt.Println("Error removing deployment: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("Removed deployment named: " + deploymentName)
}

func (c *Controller) DeploymentMongoStat(deploymentSlug string) {
	hostRegex := regexp.MustCompile(".(?:mongohq|mongolayer).com")
	loopCount := 0

	outputFormatter := func(mongoStats map[string]MongoStat, err error) {
		if err != nil {
			fmt.Println("Error authenticating connection to websocket.")
			cliOSExit()
			return
		}

		hostLength := 0
		lockLength := 0

		// Preformatting run
		for host, stats := range mongoStats {
			host = hostRegex.ReplaceAllLiteralString(host, "")
			if len(host) > hostLength {
				hostLength = len(host)
			}

			if len(stats.Locked) > lockLength {
				lockLength = len(stats.Locked) + 1
			}
		}

		headerFormat := "%" + strconv.Itoa(hostLength) + "s" + "%7s%7s%7s%7s%8s%8s%8s%8s%8s%8s%7s%" + strconv.Itoa(lockLength) + "s%11s%6s|%-3s%6s|%-3s%7s%7s%6s%11s\n"
		sprintfFormat := "%" + strconv.Itoa(hostLength) + "s" + "%7s%7s%7s%7s%8s%8s%8d%8s%8s%8s%7d%" + strconv.Itoa(lockLength) + "s%11d%6d|%-3d%6d|%-3d%7s%7s%6d%11s\n"

		if loopCount%5 == 0 {
			fmt.Printf(headerFormat, "host", "insert", "query", "update", "delete", "getmore", "command", "flush", "mapped", "vsize", "res", "faults", "locked %", "idx miss %", "qr", "qw", "ar", "aw", "netIn", "netOut", "conn", "time")
		}

		now := time.Now()

		for host, stat := range mongoStats {
			fmt.Printf(sprintfFormat, hostRegex.ReplaceAllLiteralString(host, ""), stat.Inserts, stat.Query, stat.Update, stat.Delete, stat.Getmore, stat.Command, stat.Flushes, stat.PrettyMapped(), stat.PrettyVsize(), stat.PrettyRes(), stat.Faults, stat.Locked, stat.IdxMiss, stat.Qr, stat.Qw, stat.Ar, stat.Aw, stat.PrettyNetIn(), stat.PrettyNetOut(), stat.Conn, now.Format("15:04:05"))
		}

		loopCount += 1
	}

	err := c.Api.DeploymentMongostat(deploymentSlug, outputFormatter)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		cliOSExit()
		return
	}
}

func (c *Controller) DeploymentOplog(deploymentSlug string) {
	outputFormatter := func(entry string, err error) {
		fmt.Println(entry)
	}

	err := c.Api.DeploymentOplog(deploymentSlug, outputFormatter)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		cliOSExit()
		return
	}
}
