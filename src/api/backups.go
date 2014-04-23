package api

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

func GetBackups(filter map[string]string, oauthToken string) ([]Backup, error) {
  queryString := "?"
  for key, value := range filter {
    queryString += key + "=" + value
  }
  body, err := rest_get(api_url("/backups" + queryString), oauthToken)

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

  body, err := rest_post(api_url("/backups/" + backup.Id + "/restore"), data, oauthToken)
  if err != nil {
    return Database{}, err
  }

  var database Database
  err = json.Unmarshal(body, &database)
  return database, err
}
