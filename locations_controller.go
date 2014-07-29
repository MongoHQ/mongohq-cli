package main

import (
	"fmt"
)

func (c *Controller) ListLocations() {
	providersLocations, err := c.Api.GetLocations()

	if err != nil {
		fmt.Println("Error returning locations: " + err.Error())
		cliOSExit()
		return
	}

	fmt.Println("== locations")
	for _, location := range providersLocations {
		fmt.Println("  " + location)
	}
}
