package main

import (
	"embed"
	"io/ioutil"
)

//go:embed templates
var templatesFS embed.FS

func copyTemplateFile(path, file string) error {

	// check if file exists already.
	data, err := templatesFS.ReadFile(path)
	if err != nil {
		gracefullyExit(err)
	}

	err = copyDataToFile(data, file)
	if err != nil {
		gracefullyExit(err)
	}
	return nil

}

func copyDataToFile(data []byte, toFile string) error {
	err := ioutil.WriteFile(toFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
