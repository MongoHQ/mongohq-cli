package main

import (
	"encoding/json"
)

type BackupLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type Backup struct {
	Id            string       `json:"id"`
	CreatedAt     string       `json:"created_at"`
	DatabaseNames []string     `json:"database_names"`
	DeploymentId  string       `json:"deployment_id"`
	Type          string       `json:"type"`
	Filename      string       `json:"filename"`
	Size          float64      `json:"size"`
	Links         []BackupLink `json:"links"`
	Api           Api
}

func (b *Backup) DownloadLink() string {
	for _, link := range b.Links {
		if link.Rel == "download" {
			return link.Href
		}
	}
	return "<Unknown download link>"
}

func (b *Backup) PrettySize() string {
	return prettySize(b.Size)
}

func (api *Api) GetBackups(deploymentName string) ([]Backup, error) {
	var path string

	if deploymentName != "<string>" { // this is the default returned by CLi package
		path = "/deployments/" + deploymentName + "/backups"
	} else {
		path = "/backups"
	}
	body, err := api.restGet(api.apiUrl(path))

	if err != nil {
		return []Backup{}, err
	}
	var databaseBackupSlice []Backup
	err = json.Unmarshal(body, &databaseBackupSlice)
	return databaseBackupSlice, err
}

func (api *Api) RestoreBackup(backup Backup, source, destination string) (Deployment, error) {
	type RestoreBackupParams struct {
		DestinationDatabase string `json:"name"`
		SourceDatabase      string `json:"source_name"`
	}

	restoreParams := RestoreBackupParams{DestinationDatabase: destination, SourceDatabase: source}
	data, err := json.Marshal(restoreParams)
	if err != nil {
		return Deployment{}, err
	}

	body, err := api.restPost(api.apiUrl("/backups/"+backup.Id+"/restore"), data)
	if err != nil {
		return Deployment{}, err
	}

	var deployment Deployment
	err = json.Unmarshal(body, &deployment)
	return deployment, err
}
