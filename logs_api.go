package main

import (
	"encoding/json"
	"net/url"
	"sort"
	"time"
)

type HistoricalLog struct {
	Host      string
	Message   string
	Timestamp time.Time
}

type HistoricalLogs []HistoricalLog

// Len is part of sort.Interface.
func (hl HistoricalLogs) Len() int {
	return len(hl)
}

// Swap is part of sort.Interface.
func (hl HistoricalLogs) Swap(i, j int) {
	hl[i], hl[j] = hl[j], hl[i]
}

// Less is part of sort.Interface. We use count as the value to sort by
func (hl HistoricalLogs) Less(i, j int) bool {
	return hl[i].Timestamp.Before(hl[j].Timestamp)
}

func (api *Api) GetHistoricalLogs(deploymentSlug, search, exclude, regexp string) (HistoricalLogs, int, error) {
	var historicalLogs HistoricalLogs
	maxHostnameLength := 0

	urlPath := "/deployments/" + api.Config.AccountSlug + "/" + deploymentSlug + "/historical_logs?size=200&sort=desc"

	if regexp != "<string>" {
		urlPath = urlPath + "&grep=" + url.QueryEscape(regexp)
	}

	if search != "<string>" {
		urlPath = urlPath + "&search=" + url.QueryEscape(search)
	}

	if exclude != "<string>" {
		urlPath = urlPath + "&exclude=" + url.QueryEscape(exclude)
	}

	body, err := api.restGet(api.apiUrl(urlPath))
	if err != nil {
		return nil, maxHostnameLength, err
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, maxHostnameLength, err
	}

	for host, logs := range result {
		for _, log := range logs.(map[string]interface{})["logs"].([]interface{}) {
			ts := log.(map[string]interface{})["ts"].(string)
			timestamp, _ := time.Parse("2006-01-02T15:04:05Z", ts)
			if maxHostnameLength < len(formatHostname(host)) {
				maxHostnameLength = len(formatHostname(host))
			}
			historicalLogs = append(historicalLogs, HistoricalLog{Host: host, Message: log.(map[string]interface{})["message"].(string), Timestamp: timestamp})
		}
	}
	sort.Sort(historicalLogs)
	return historicalLogs, maxHostnameLength, err
}
