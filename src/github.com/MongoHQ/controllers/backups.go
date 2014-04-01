package controllers

import (
  "github.com/MongoHQ/api"
  "fmt"
  "os"
  "strings"
  "errors"
)

func Backups(filter map[string]string) {
  backupsSlice, err := api.GetBackups(filter, OauthToken)

  if err != nil {
    fmt.Println("Error retreiving backups: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("== Backups")
  for _, backup := range backupsSlice {
    fmt.Println(backup.Filename)
  }
}

func Backup(filename string) {
  backup, err := findBackupByFilename(filename)
  if err != nil {
    fmt.Println("Error retreiving backup: " + err.Error())
    os.Exit(1)
  }
  fmt.Println("== Backup " + filename)
  fmt.Println(" id            : " + backup.Id)
  fmt.Println(" created at    : " + backup.CreatedAt)
  fmt.Println(" databases     : " + strings.Join(backup.DatabaseNames, ", "))
  fmt.Println(" type          : " + backup.Type)
  fmt.Println(" size          : " + backup.PrettySize())
}

func RestoreBackup(filename, source, destination string) {
  backup, err := findBackupByFilename(filename)
  if err != nil {
    fmt.Println("Error retreiving backup: " + err.Error())
    os.Exit(1)
  }

  database, err := backup.Restore(source, destination, OauthToken)
  if err != nil {
    fmt.Println("Error restoring backup: " + err.Error())
    os.Exit(1)
  }

  fmt.Println("=== Restoring from database " + source + " in backup file " + filename + " to new database " + destination)
  pollNewDeployment(database)
}

func findBackupByFilename(filename string) (api.Backup, error) {
  backups, err := api.GetBackups(map[string]string{}, OauthToken)
  if err != nil {
    return api.Backup{}, err
  }
  for _, backup := range backups {
    if backup.Filename == filename {
      return backup, err
    }
  }
  return api.Backup{}, errors.New("Could not find backup with value " + filename)
}
