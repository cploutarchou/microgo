package main

import (
	"errors"
	"github.com/fatih/color"
	"os"

	"github.com/cploutarchou/microGo"
)

const version = "1.0.0"

var micro microGo.MicroGo

func main() {
	var message string
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		gracefullyExit(err)
	}

	setup()

	switch arg1 {
	case "help":
		help()
	case "version":
		color.Yellow("Application version: " + version)

	case "migrate":
		if arg2 == "" {
			arg2 = "up"
		}
		err = doMigrate(arg2, arg3)
		if err != nil {
			gracefullyExit(err)
		}
		message = "Migration successfully completed."

	case "make":
		if arg2 == "" {
			gracefullyExit(errors.New("make command requires an argument . Available options: migration|model|handler "))
		}
		err = makeDo(arg2, arg3)
		if err != nil {
			gracefullyExit(err)
		}
	default:
		help()
	}
	gracefullyExit(nil, message)
}

func validateInput() (string, string, string, error) {
	var arg1, arg2, arg3 string
	if len(os.Args) > 1 {
		arg1 = os.Args[1]
		if len(os.Args) >= 3 {
			arg2 = os.Args[2]
		}
		if len(os.Args) >= 4 {
			arg3 = os.Args[3]
		}
	} else {
		color.Red("Error : No valid input.")
		help()
		return "", "", "", errors.New("Command line arguments is required! ")
	}
	return arg1, arg2, arg3, nil
}

func gracefullyExit(err error, msg ...string) {
	message := ""
	if len(msg) > 0 {
		message = msg[0]
	}

	if err != nil {
		color.Red("Error: %v\n", err)
	}
	if len(message) > 0 {
		color.Yellow(message)
	} else {
		color.Green("Completed")
	}
	os.Exit(0)
}
