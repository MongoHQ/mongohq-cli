package main

import (
	"fmt"
	"os"
)

func (c *Controller) CurrentUser() {
	user, err := c.Api.GetCurrentUser()

	if err != nil {
		fmt.Println("Error returning user: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("== whoami")
	fmt.Println("  name:  " + user.Name)
	fmt.Println("  email: " + user.Email)
}
