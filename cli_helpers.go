package main

import (
  "fmt"
  "os"
	"github.com/codegangsta/cli"
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
  fmt.Println(command)
}
