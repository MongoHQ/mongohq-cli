package main

import (
	"fmt"
	"os"
)

func (c *Controller) SetConfigAccount(slug string) {
	account, err := c.Api.GetAccount(slug)

	if err != nil {
		fmt.Println("Error accessing account:" + err.Error())
		cliOSExit()
		return
	}

	config := getConfig()
	config.AccountSlug = account.Slug

	if err := config.Save(); err != nil {
		fmt.Println("Error setting default account: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("Set default account to " + account.Slug)
}

func requireAccount(api *Api) {
	runCount := 0
	for api.Config.AccountSlug == "" {
		if runCount > 2 {
			fmt.Println("Default account is required.  Please run `mongohq accounts` and `mongohq config:account -a <account-slug>` to set a default acount.")
			os.Exit(1)
		}

		accounts, err := api.GetAccounts()

		if err != nil {
			fmt.Println("Error returning list of accounts let's try one more time.")
			runCount += 1
			continue
		}

		var account Account

		if len(accounts) > 1 {
			fmt.Println("To continue, we need to set a default account.  Here is a list of accounts you have access to:")
			for _, account := range accounts {
				fmt.Println("  " + account.Slug)
			}

			accountSlug := prompt("Which account should be default")
			account, err = api.GetAccount(accountSlug)

			if err != nil {
				fmt.Println("Error accessing the account '" + accountSlug + "'.  Please try again.")
				runCount += 1
				continue
			}

		} else {
			account = accounts[0]
		}

		config := getConfig()
		config.AccountSlug = account.Slug
		err = config.Save()

		api.Config = config

		if err != nil {
			fmt.Println("Error saving default configuartion to ~/.mongohq/defaults. Will you check the permissions on the file allow you to write?")
		}

		fmt.Println("Set default account to " + account.Slug + "\n")
		return
	}
}
