package main

import (
	"encoding/json"
	"strings"
)

type Deployment struct {
	Id                     string
	Name                   string `json:"name"`
	Plan                   string `json:"plan"`
	Location               string `json:"location"`
	CurrentPrimary         string `json:"current_primary"`
	Status                 string `json:"status"`
	Version                string
	Members                []string
	AllowMultipleDatabases bool `json:"allow_multiple_deployments"`
	Databases              []Database
}

func (d *Deployment) NameOrId() string {
	if d.Name != "" {
		return d.Name
	} else {
		return d.Id
	}
}

type SocketMessage struct {
	Command string  `json:"command"`
	Uuid    string  `json:"uuid"`
	Message Message `json:"message"`
}

type Message struct {
	Deployment string `json:"deployment"`
	Account    string `json:"account"`
	Type       string `json:"type"`
}

type MongoStatMessage struct {
	Type    string
	Ts      string
	Error   string
	Message map[string]MongoStat
}

type MongoStat struct {
	Version    string
	Inserts    string
	RawInserts int `json:"raw_inserts"`
	Query      string
	RawQuery   int `json:"raw_query"`
	Update     string
	RawUpdate  int `json:"raw_update"`
	Delete     string
	RawDelete  int `json:"raw_delete"`
	Getmore    string
	RawGetmore int `json:"raw_getmore"`
	Command    string
	RawCommand int `json:"raw_command"`
	Flushes    int
	Mapped     int
	Vsize      int
	Res        int
	Faults     int
	Locked     string
	IdxMiss    int `json:"idx_miss"`
	Qr         int
	Qw         int
	Ar         int
	Aw         int
	NetIn      float64 `json:"net_in"`
	NetOut     float64 `json:"net_out"`
	Conn       int
	Set        string
	Repl       string
}

func (m *MongoStat) PrettyNetIn() string {
	return prettySize(float64(m.NetIn))
}

func (m *MongoStat) PrettyNetOut() string {
	return prettySize(float64(m.NetOut))
}

func (m *MongoStat) PrettyMapped() string {
	return prettySize(float64(m.Mapped * 1024 * 1024))
}

func (m *MongoStat) PrettyVsize() string {
	return prettySize(float64(m.Vsize * 1024 * 1024))
}

func (m *MongoStat) PrettyRes() string {
	return prettySize(float64(m.Res * 1024 * 1024))
}

func (api *Api) GetDeployments() ([]Deployment, error) {
	body, err := api.restGet(api.apiUrl("/accounts/" + api.Config.AccountSlug + "/deployments"))

	if err != nil {
		return nil, err
	}
	var deploymentsSlice []Deployment
	err = json.Unmarshal(body, &deploymentsSlice)
	return deploymentsSlice, err
}

func (api *Api) GetDeployment(deploymentId string) (Deployment, error) {
	body, err := api.restGet(api.apiUrl("/deployments/"+api.Config.AccountSlug+"/"+deploymentId) + "?embed=databases,plan")

	if err != nil {
		return Deployment{}, err
	}
	var deployment Deployment
	err = json.Unmarshal(body, &deployment)
	return deployment, err
}

func (api *Api) CreateDeployment(deploymentName, databaseName, location string) (Deployment, error) {
	type DeploymentCreate struct {
		Name         string `json:"name"`
		DatabaseName string `json:"database_name"`
		Location     string `json:"location"`
		Type         string `json:"type"`
	}

	deploymentCreate := DeploymentCreate{Name: deploymentName, DatabaseName: databaseName, Location: location, Type: "mongodb"}
	data, err := json.Marshal(deploymentCreate)
	if err != nil {
		return Deployment{}, err
	}

	body, err := api.restPost(api.apiUrl("/accounts/"+api.Config.AccountSlug+"/deployments/elastic"), data)

	if err != nil {
		return Deployment{}, err
	}
	var deployment Deployment
	err = json.Unmarshal(body, &deployment)
	return deployment, err
}

func (api *Api) RemoveDeployment(deploymentSlug string) error {
	_, err := api.restDelete(api.apiUrl("/deployments/" + api.Config.AccountSlug + "/" + deploymentSlug))
	return err
}

func (api *Api) RenameDeployment(deploymentId, name string) (Deployment, error) {
	type DeploymentRenameParams struct {
		Name string `json:"name"`
	}

	data, err := json.Marshal(DeploymentRenameParams{Name: name})
	if err != nil {
		return Deployment{}, err
	}

	body, err := api.restPatch(api.apiUrl("/deployments/"+api.Config.AccountSlug+"/"+deploymentId), data)
	if err != nil {
		return Deployment{}, err
	}
	var deployment Deployment
	err = json.Unmarshal(body, &deployment)
	return deployment, err
}

func (api *Api) BackupDeployment(deploymentId string) (Backup, error) {
	var err error
	body, err := api.restPost(api.apiUrl("/deployments/"+api.Config.AccountSlug+"/"+deploymentId+"/backups"), []byte{})
	if err != nil {
		return Backup{}, err
	}
	var backup Backup
	err = json.Unmarshal(body, &backup)
	return backup, err
}

func (api *Api) DeploymentMongostat(deploymentSlug string, outputFormatter func(map[string]MongoStat, error)) error {
	message := SocketMessage{Command: "subscribe", Uuid: "12345", Message: Message{Account: api.Config.AccountSlug, Deployment: deploymentSlug, Type: "mongo.stats"}}
	socket, err := api.openWebsocket(message)
	if err != nil {
		return err
	}

	for {
		_, msg, err := socket.ReadMessage()
		if err != nil {
			outputFormatter(make(map[string]MongoStat, 0), err)
		}

		// catch the first success response
		if strings.Index(string(msg), "successful") > -1 {
			continue
		}

		// null is bad news for Go, and gopher has an outstanding issue with
		// the first response: https://github.com/MongoHQ/gopher/issues/14
		if strings.Index(string(msg), "null") > -1 {
			continue
		}

		mongoStatMessage := MongoStatMessage{}
		err = json.Unmarshal(msg, &mongoStatMessage)
		outputFormatter(mongoStatMessage.Message, err)
	}
}

func (api *Api) DeploymentOplog(deploymentSlug string, outputFormatter func(string, error)) error {
	message := SocketMessage{Command: "subscribe", Uuid: "12345", Message: Message{Deployment: deploymentSlug, Type: "mongo.oplog"}}
	socket, err := api.openWebsocket(message)
	if err != nil {
		return err
	}

	for {
		_, msg, err := socket.ReadMessage()
		outputFormatter(string(msg), err)
	}
}
