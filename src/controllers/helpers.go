package controllers

import (
	"fmt"
	"regexp"
)

func prompt(text string) string {
	var response string
	var err error

	print(text + ": ")
	_, err = fmt.Scanln(&response)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return response
}

func formatHostname(host string) string {
	hostRegex := regexp.MustCompile(".(?:mongohq|mongolayer).com")
	host = hostRegex.ReplaceAllLiteralString(host, "")
	return host
}
