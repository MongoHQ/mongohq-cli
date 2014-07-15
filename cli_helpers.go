package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
)

func requireArguments(c *cli.Context, argumentsSlice []string, errorMessages []string) {
	err := false

	for _, argument := range argumentsSlice {
		if !c.IsSet(argument) {
			err = true
			fmt.Println("--" + argument + " is required")
		}
	}

	if err {
		fmt.Println("\nMissing arguments, for more information, run: mongohq " + c.Command.Name + " --help\n")
		for _, errorMessage := range errorMessages {
			fmt.Println(errorMessage)
		}
		os.Exit(1)
	}
}

func findClosestCommand(context *cli.Context, command string) {
	fmt.Println(" ! `" + command + "` is not a mongohq command.")
	fmt.Println(" ! See `mongohq help` for a list of available commands")
	os.Exit(1)
}
