package main

import (
	"fmt"
	"os"
)

func (c *Controller) ListLocations() {
	providersLocations, err := c.Api.GetLocations()

	if err != nil {
		fmt.Println("Error returning locations: " + err.Error())
		os.Exit(1)
	}

	fmt.Println("== locations")
	for _, location := range providersLocations {
		fmt.Println("  " + location)
	}
}
