package main

import (
  "fmt"
  "os"
  "strings"
  "errors"
)

func (c *Controller) ListBackups(filter map[string]string) {
  backupsSlice, err := c.Api.GetBackups(filter, OauthToken)

  if err != nil {
    fmt.Println("Error retreiving backups: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("== Backups")
  for _, backup := range backupsSlice {
    fmt.Println(backup.Filename)
  }
}

func (c *Controller) ShowBackup(filename string) {
  backup, err := c.findBackupByFilename(filename)
  if err != nil {
    fmt.Println("Error retreiving backup: " + err.Error())
    os.Exit(1)
  }
  fmt.Println("== Backup " + filename)
  fmt.Println(" created at    : " + backup.CreatedAt)
  fmt.Println(" databases     : " + strings.Join(backup.DatabaseNames, ", "))
  fmt.Println(" type          : " + backup.Type)
  fmt.Println(" size          : " + backup.PrettySize())
}

func (c *Controller) RestoreBackup(filename, source, destination string) {
  backup, err := c.findBackupByFilename(filename)
  if err != nil {
    fmt.Println("Error retreiving backup: " + err.Error())
    os.Exit(1)
  }

  database, err := c.Api.RestoreBackup(backup, source, destination)
  if err != nil {
    fmt.Println("Error restoring backup: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("=== Restoring from database " + source + " in backup file " + filename + " to new database " + destination)
  c.pollNewDeployment(database)
}

func (c *Controller) findBackupByFilename(filename string) (Backup, error) {
  backups, err := c.Api.GetBackups(map[string]string{}, OauthToken)
  if err != nil {
    return Backup{}, err
  }
  for _, backup := range backups {
    if backup.Filename == filename {
      return backup, err
    }
  }
  return Backup{}, errors.New("Could not find backup with value " + filename)
}
