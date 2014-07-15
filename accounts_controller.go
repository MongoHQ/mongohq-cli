package main

import (
	"fmt"
	"os"
)

func (c *Controller) ListAccounts() {
	accountsSlice, err := c.Api.GetAccounts()

	if err != nil {
		fmt.Println("Error retreiving accounts: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("== Accounts")
	for _, account := range accountsSlice {
		if c.Api.Config.AccountSlug == account.Slug { // signify it is the default account
			fmt.Println(account.Slug + " (default)")
		} else {
			fmt.Println(account.Slug)
		}
	}
}

func (c *Controller) ShowAccount(slug string) {
	account, err := c.Api.GetAccount(slug)

	if err != nil {
		fmt.Println("Error retreiving account: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("== " + slug)
	fmt.Println(" slug:    " + account.Slug)
	fmt.Println(" name:    " + account.Name)
	fmt.Println(" owner:   " + account.OwnerName + " (" + account.OwnerEmail + ") ")
	fmt.Println(" == Users")
	for _, user := range account.Users {
		fmt.Println("    " + user.Name + " (" + user.Email + ") ")
	}
}
