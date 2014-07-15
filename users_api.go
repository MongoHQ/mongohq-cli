package main

import (
	"encoding/json"
)

type User struct {
	Id       string    `json:"id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Accounts []Account `json:"accounts"`
}

func (api *Api) GetCurrentUser() (User, error) {
	body, err := api.restGet(api.apiUrl("/user"))

	if err != nil {
		return User{}, err
	}

	var user User
	err = json.Unmarshal(body, &user)
	return user, err
}
