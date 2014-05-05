package main

import (
  "encoding/json"
)

type BackupLink struct {
  Rel string `json:"rel"`
  Href string `json:"href"`
}

type Backup struct {
  Id string `json:"id"`
  CreatedAt string `json:"created_at"`
  DatabaseNames []string `json:"database_names"`
  Type string `json:"type"`
  Filename string `json:"filename"`
  Size float64 `json:"size"`
  Links []BackupLink `json:"links"`
  Api Api
}

func (b *Backup) DownloadLink() string {
  for _, link := range b.Links {
    if link.Rel == "download" { return link.Href }
  }
  return "<Unknown download link>"
}

func (b *Backup) PrettySize() string {
  return prettySize(b.Size)
}

func (api *Api) GetBackups(filter map[string]string, oauthToken string) ([]Backup, error) {
  queryString := "?"
  for key, value := range filter {
    queryString += key + "=" + value
  }
  body, err := api.restGet(api.apiUrl("/backups" + queryString))

  if err != nil {
    return []Backup{}, err
  }
  var databaseBackupSlice []Backup
  err = json.Unmarshal(body, &databaseBackupSlice)
  return databaseBackupSlice, err
}

func (backup *Backup) Restore(source, destination, oauthToken string) (Database, error) {
  type RestoreBackupParams struct {
    DestinationDatabase string `json:"name"`
    SourceDatabase string `json:"source_name"`
  }

  restoreParams := RestoreBackupParams{DestinationDatabase: destination, SourceDatabase: source}
  data, err := json.Marshal(restoreParams)
  if err != nil {
    return Database{}, err
  }

  body, err := backup.Api.restPost(backup.Api.apiUrl("/backups/" + backup.Id + "/restore"), data)
  if err != nil {
    return Database{}, err
  }

  var database Database
  err = json.Unmarshal(body, &database)
  return database, err
}
