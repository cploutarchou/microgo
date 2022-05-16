package main

import (
	"embed"
	"errors"
	"io/ioutil"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyTemplateFile(templatePath, targetFile string) error {
	if fileExists(targetFile) {
		return errors.New(targetFile + " already exists!")
	}

	data, err := templateFS.ReadFile(templatePath)
	if err != nil {
		gracefullyExit(err)
	}

	err = copyDataToFile(data, targetFile)
	if err != nil {
		gracefullyExit(err)
	}

	return nil
}

func copyDataToFile(data []byte, to string) error {
	err := ioutil.WriteFile(to, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}
