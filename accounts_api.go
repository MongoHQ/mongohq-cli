package main

import (
	"encoding/json"
)

type Account struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	OwnerId   string `json:"owner_id"`
	OwnerName string `json:"owner"`
}

func (api *Api) GetAccounts() ([]Account, error) {
	body, err := api.restGet(api.apiUrl("/accounts"))

	if err != nil {
		return []Account{}, err
	}
	var accountsSlice []Account
	err = json.Unmarshal(body, &accountsSlice)
	return accountsSlice, err
}

func (api *Api) GetAccount(slug string) (Account, error) {
	body, err := api.restGet(api.apiUrl("/" + slug + "/settings"))

	if err != nil {
		return Account{}, err
	}
	var account Account
	err = json.Unmarshal(body, &account)
	return account, err
}
