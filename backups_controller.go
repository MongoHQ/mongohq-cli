package main

import (
	"fmt"
	"strings"
)

func (c *Controller) ListBackups() {
	backupsSlice, err := c.Api.GetBackups()

	if err != nil {
		fmt.Println("Error retreiving backups: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== Backups")
	for _, backup := range backupsSlice {
		fmt.Println(backup.Filename)
	}
}

func (c *Controller) ListBackupsForDeployment(deploymentSlug string) {
	backupsSlice, err := c.Api.GetBackupsForDeployment(deploymentSlug)

	if err != nil {
		fmt.Println("Error retreiving backups: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== Backups for " + deploymentSlug)
	for _, backup := range backupsSlice {
		fmt.Println(backup.Filename)
	}
}

func (c *Controller) ShowBackup(backupSlug string) {
	backup, err := c.Api.GetBackup(backupSlug)
	if err != nil {
		fmt.Println("Error retreiving backup: " + err.Error())
		cliOSExit()
		return
	}
	deployment, _ := c.Api.GetDeployment(backup.DeploymentSlug)
	fmt.Println("== Backup " + backupSlug)
	fmt.Println(" deployment : " + deployment.Name)
	fmt.Println(" databases  : " + strings.Join(backup.DatabaseNames, ", "))
	fmt.Println(" status     : " + backup.Status)
	fmt.Println(" created at : " + backup.CreatedAt)
	fmt.Println(" type       : " + backup.Type)
	fmt.Println(" size       : " + backup.PrettySize())
	if backup.Status == "complete" {
		fmt.Println(" download   : " + backup.DownloadLink())
	}
}

func (c *Controller) RestoreBackup(backupSlug, deploymentName, source, destination string) {
	backup, err := c.Api.GetBackup(backupSlug)
	if err != nil {
		fmt.Println("Error retreiving backup: " + err.Error())
		cliOSExit()
		return
	}

	deployment, err := c.Api.RestoreBackup(backup, deploymentName, source, destination)
	if err != nil {
		fmt.Println("Error restoring backup: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== Restoring from database " + source + " on deployment " + backup.DeploymentSlug + " from backup " + backupSlug + " to new deployment " + deploymentName)
	c.pollNewDeployment(deployment)
}

func (c *Controller) CreateBackup(deploymentSlug string) {
	backup, err := c.Api.BackupDeployment(deploymentSlug)
	if err != nil {
		fmt.Println("Error triggering backup on deployment: " + err.Error())
		cliOSExit()
		return
	}

	status := backup.Status
	fmt.Print("Running backup")
	for status == "running" {
		fmt.Print(".")
		backup, err = c.Api.GetBackup(backup.Id)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("\nError requesting backup status. For a manual update, please run:\n\n compose backups:info -b " + backup.Id)
			cliOSExit()
			return
		}
		status = backup.Status
	}

	if status != "complete" {
		fmt.Println("Error creating backup.  Please try once more, or contact support@compose.io.")
		cliOSExit()
		return
	}

	fmt.Println("== Backup " + backup.Filename)
	fmt.Println(" deployment : " + deploymentSlug)
	fmt.Println(" databases  : " + strings.Join(backup.DatabaseNames, ", "))
	fmt.Println(" status     : " + backup.Status)
	fmt.Println(" created at : " + backup.CreatedAt)
	fmt.Println(" type       : " + backup.Type)
	fmt.Println(" size       : " + backup.PrettySize())
	if backup.Status == "complete" {
		fmt.Println(" download   : " + backup.DownloadLink())
	}
}
