package main

import (
	"fmt"
	"github.com/codegangsta/cli"
)

func requireArguments(c *cli.Context, argumentsSlice []string, errorMessages []string) error {
	err := false

	for _, argument := range argumentsSlice {
		if !c.IsSet(argument) {
			err = true
			fmt.Println("--" + argument + " is required")
		}
	}

	if err {
		if !replMode {
			fmt.Println("\nMissing arguments, for more information, run: mongohq " + c.Command.Name + " --help\n")
		} else {
			fmt.Println("Missing arguments: type 'help " + c.Command.Name + "' for details")
		}

		for _, errorMessage := range errorMessages {
			fmt.Println(errorMessage)
		}

		return fmt.Errorf("Missing arguments")
	}
	return nil
}

func findClosestCommand(context *cli.Context, command string) {

	if !replMode {
		fmt.Println(" ! `" + command + "` is not a mongohq command.")
		fmt.Println(" ! See `mongohq help` for a list of available commands")
	} else {
		fmt.Println("Unknown command:" + command)
	}
	cliOSExit()
}
