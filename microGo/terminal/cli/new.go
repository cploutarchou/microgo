package main

import (
	"errors"
	"strings"
)

func createNew(applicationName string) {
	applicationName = strings.TrimSpace(applicationName)
	applicationName = strings.ToLower(applicationName)
	if applicationName == "" {
		gracefullyExit(errors.New("No project name specified! "))
	}

	// sanitize the application name

	// git clone skeleton application

	//remove the .git directory

	// create a new .env file

	// create a makefile for the application

	// update the go.mod file

	// update the existing .go files with th correct package names

	// run go mod tidy

}
