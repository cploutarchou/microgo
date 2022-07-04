package main

import (
	"errors"
	"log"
	"strings"
)

func createNew(applicationName string) {
	applicationName = strings.TrimSpace(applicationName)
	applicationName = strings.ToLower(applicationName)
	if applicationName == "" {
		gracefullyExit(errors.New("No project name specified! "))
	}

	// sanitize the application name
	if strings.Contains(applicationName, "/") {
		exploded := strings.SplitAfter(applicationName, "/")
		applicationName = exploded[len(exploded)-1]
	}
	log.Println("Application name: ", applicationName)
	// git clone skeleton application

	//remove the .git directory

	// create a new .env file

	// create a makefile for the application

	// update the go.mod file

	// update the existing .go files with th correct package names

	// run go mod tidy

}
