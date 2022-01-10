package main

import (
	"errors"
	"github.com/fatih/color"
	"log"
	"os"
)

const version = "1.0.0"

//var micro microGo.MicroGo

func main() {
	arg1, arg2, arg3, err := validateInput()
	if err != nil {
		gracefullyExit(err)
	}
	switch arg1 {
	case "help":
		help()
	case "version":
		color.HiWhite("Application version: " + version)
	default:
		log.Println(arg2, arg3)
	}
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

func help() {
	color.HiWhite(`Available commands:`)
	color.Yellow(`help             - Shows the help commands
version          - Print application version`,
	)
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
