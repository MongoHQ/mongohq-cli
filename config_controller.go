package main

import (
	"fmt"
	"os"
)

func (c *Controller) SetConfigAccount(slug string) {
	account, err := c.Api.GetAccount(slug)

	if err != nil {
		fmt.Println("Error retreiving account: " + err.Error())
		os.Exit(1)
	}

	config := getConfig()
	config.AccountSlug = account.Slug

	if err := config.Save(); err != nil {
		fmt.Println("Error setting default account: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Set default account to " + account.Slug)
}

func requireAccount(c Config) {
	if c.AccountSlug == "" {
		fmt.Println("Default account is required.  Please run `mongohq accounts` and `mongohq config:account -a <account-slug>` to set a default acount.")
		os.Exit(1)
	}
}
