package main

import (
	"code.google.com/p/gopass"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/peterh/liner"
	"os"
	"os/signal"
	"strings"
)

var replMode bool // This is the repl mode - set to avoid exits
var replExit bool // This is the flag to set if you need the repl to exit
var term *liner.State
var myapp *cli.App

func repl(app *cli.App) {
	myapp = app
	replMode = true
	replExit = false
	prompt := "> "

	initTerm()

	for !replExit {
		line, err := term.Prompt(prompt)
		if err != nil {
			break
		}
		if line == "exit" {
			break
		}
		term.AppendHistory(line)
		plines := strings.Fields(line)
		plines = append(plines, "")
		copy(plines[1:], plines[0:])
		plines[0] = "mongohq-repl"
		app.Run(plines)
	}

}

func safeGetPass(prompt string) (passwd string, err error) {

	if replMode {
		closeTerm()
	}
	passwd, err = gopass.GetPass(prompt)

	if replMode {
		initTerm()
	}

	return passwd, err
}

func initTerm() {
	term = liner.NewLiner()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		closeTerm()
		fmt.Println("Quitting shell")
		os.Exit(1)
	}()

	if f, err := os.Open(historyfn); err == nil {
		term.ReadHistory(f)
		f.Close()
	}

	// TODO: This is a basic completer for *only* the subcommands
	term.SetCompleter(func(line string) (c []string) {
		// TODO:
		// split line, count args, if none, sub command match
		// if more than one, check what the last entered thing is
		// if -- load up with the flags
		// TODO++: otherwise look at previous flag and determine
		// possible completetions
		parts := strings.Fields(line)
		switch {
		case len(parts) == 0:
			// Nothing - load all the lines
			for _, cmd := range myapp.Commands {
				c = append(c, cmd.Name)
			}

		case len(parts) == 1:
			// One item - if it is help then offer help for all commands
			if parts[0] == "help" {
				for _, cmd := range myapp.Commands {
					c = append(c, "help "+cmd.Name)
				}
			} else {
				for _, cmd := range myapp.Commands {
					if strings.HasPrefix(cmd.Name, strings.ToLower(line)) {
						c = append(c, cmd.Name)
					}
				}
			}

		case len(parts) == 2:
			// two parts, so is this help command-partial if so
			if parts[0] == "help" {
				for _, cmd := range myapp.Commands {
					if strings.HasPrefix(cmd.Name, strings.ToLower(parts[1])) {
						c = append(c, "help "+cmd.Name)
					}
				}
			}
		}
		return
	})
}

func closeTerm() {
	if f, err := os.Create(historyfn); err != nil {
		fmt.Println("Error writing history file", err)
	} else {
		term.WriteHistory(f)
		f.Close()
	}
	term.Close()
}

func cliOSExit() {
	if !replMode {
		os.Exit(1)
	}
}
