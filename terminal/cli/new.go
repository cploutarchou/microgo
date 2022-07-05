package main

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"io"
	"log"
	"os"
	"runtime"
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
	color.Green("\tCloning skeleton application from git repository...")

	_, err := git.PlainClone("./"+applicationName, false, &git.CloneOptions{
		URL:      "https://github.com/cploutarchou/microGo_skeleton_app.git",
		Progress: os.Stdout,
		Depth:    1,
	})
	if err != nil {
		gracefullyExit(err)
	}
	//remove the .git directory
	err = os.RemoveAll(fmt.Sprintf("./%s/.git", applicationName))
	if err != nil {
		gracefullyExit(err)
	}
	// create a new .env file
	color.Yellow("Creating a new .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		gracefullyExit(err)
	}
	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", applicationName)
	env = strings.ReplaceAll(env, "${KEY}", micro.CreateRandomString(32))
	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", applicationName))
	if err != nil {
		gracefullyExit(err)
	}

	// create a makefile for the application
	if runtime.GOOS == "windows" {
		src, err := os.Open(fmt.Sprintf("./%s/Makefile.windows", applicationName))
		if err != nil {
			gracefullyExit(err)
		}
		defer func(src *os.File) {
			err := src.Close()
			if err != nil {
				gracefullyExit(err)
			}
		}(src)
		dest, err := os.Create(fmt.Sprintf("./%s/Makefile", applicationName))
		if err != nil {
			gracefullyExit(err)
		}
		defer func(dest *os.File) {
			err := dest.Close()
			if err != nil {

			}
		}(dest)
		_, err = io.Copy(dest, src)
		if err != nil {
			gracefullyExit(err)
		}
	} else {
		src, err := os.Open(fmt.Sprintf("./%s/Makefile.mac", applicationName))
		if err != nil {
			gracefullyExit(err)
		}
		defer func(src *os.File) {
			err := src.Close()
			if err != nil {

			}
		}(src)
		dest, err := os.Create(fmt.Sprintf("./%s/Makefile", applicationName))
		if err != nil {
			gracefullyExit(err)
		}
		defer func(dest *os.File) {
			err := dest.Close()
			if err != nil {

			}
		}(dest)
		_, err = io.Copy(dest, src)
		if err != nil {
			gracefullyExit(err)
		}

	}
	_ = os.Remove("./" + applicationName + "/Makefile.mac")
	_ = os.Remove("./" + applicationName + "/Makefile.windows")

	// update the go.mod file

	// update the existing .go files with th correct package names

	// run go mod tidy

}
