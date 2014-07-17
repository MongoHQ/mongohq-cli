package main

import (
	"encoding/json"
)

type Database struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Plan         string `json:"plan"`
	DeploymentId string `json:"deployment_id"`
}

type DatabaseUser struct {
	Username     string `json:"user"`
	PasswordHash string `json:"pwd"`
	ReadOnly     bool
}

type DatabaseStats struct {
	Database          string         `json:"db"`
	AverageObjectSize float64        `json:"avgObjSize"`
	Collections       int            `json:"collections"`
	DataFileVersion   map[string]int `json:"dataFileVersion"`
	DataSize          int            `json:"dataSize"`
	FileSize          int            `json:"fileSize"`
	IndexSize         int            `json:"indexSize"`
	Indexes           int            `json:"indexes"`
	NsSizeMb          int            `json:"nsSizeMB"`
	NumExtents        int            `json:"numExtents"`
	Objects           int            `json:"objects"`
	Ok                int            `json:"ok"`
	StorageSize       int            `json:"storageSize"`
}

func (api *Api) GetDatabases() ([]Database, error) {
	body, err := api.restGet(api.apiUrl("/databases"))

	if err != nil {
		return nil, err
	}
	var databasesSlice []Database
	err = json.Unmarshal(body, &databasesSlice)
	return databasesSlice, err
}

func (api *Api) GetDatabase(deploymentName, databaseName string) (Database, error) {
	body, err := api.restGet(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + deploymentName + "/databases/" + databaseName))

	if err != nil {
		return Database{}, err
	}
	var database Database
	err = json.Unmarshal(body, &database)
	return database, err
}

func (api *Api) GetDatabaseStats(database Database) (map[string]DatabaseStats, error) {
	body, err := api.restGet(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + database.DeploymentId + "/mongodb/" + database.Name + "/stats"))

	if err != nil {
		return make(map[string]DatabaseStats), err
	}

	var dbStats map[string]DatabaseStats
	err = json.Unmarshal(body, &dbStats)
	return dbStats, err
}

func (api *Api) CreateDatabase(deploymentName, databaseName string) (Database, error) {
	type DatabaseCreate struct {
		Name string `json:"name"`
	}

	databaseCreate := DatabaseCreate{Name: databaseName}
	data, err := json.Marshal(databaseCreate)
	if err != nil {
		return Database{}, err
	}

	body, err := api.restPost(api.apiUrl("/deployments/"+api.Config.AccountSlug+"/"+deploymentName+"/databases"), data)

	if err != nil {
		return Database{}, err
	}
	var database Database
	err = json.Unmarshal(body, &database)
	return database, err
}

func (api *Api) RemoveDatabase(deploymentSlug, databaseName string) error {
	_, err := api.restDelete(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + deploymentSlug + "/databases/" + databaseName))
	return err
}

func (api *Api) GetDatabaseUsers(deployment_id, database_name string) ([]DatabaseUser, error) {
	body, err := api.restGet(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + deployment_id + "/mongodb/" + database_name + "/users"))
	if err != nil {
		return make([]DatabaseUser, 0), err
	}
	var databaseUsersSlice []DatabaseUser
	err = json.Unmarshal(body, &databaseUsersSlice)
	return databaseUsersSlice, err
}

func (api *Api) CreateDatabaseUser(deploymentId, databaseName, username, password string) (OkResponse, error) {
	type UserCreate struct {
		Username string `json:"username"`
		Password string `json:"password"`
		ReadOnly bool   `json:"readOnly"`
	}

	userCreate := UserCreate{Username: username, Password: password, ReadOnly: false}
	data, err := json.Marshal(userCreate)
	if err != nil {
		return OkResponse{}, err
	}

	body, err := api.restPost(api.apiUrl("/deployments/"+api.Config.AccountSlug+"/"+deploymentId+"/mongodb/"+databaseName+"/users"), data)

	if err != nil {
		return OkResponse{}, err
	}
	var okResponse OkResponse
	err = json.Unmarshal(body, &okResponse)
	return okResponse, err
}

func (api *Api) RemoveDatabaseUser(deploymentId, databaseName, username string) (OkResponse, error) {
	body, err := api.restDelete(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + deploymentId + "/mongodb/" + databaseName + "/users/" + username))

	if err != nil {
		return OkResponse{}, err
	}
	var okResponse OkResponse
	err = json.Unmarshal(body, &okResponse)
	return okResponse, err
}
