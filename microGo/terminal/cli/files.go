package main

import (
	"embed"
	"io/ioutil"
)

//go:embed templates
var templatesFS embed.FS

func copyTemplateFile(templatesPath, targetFile string) error {

	data, err := templatesFS.ReadFile(templatesPath)
	if err != nil {
		gracefullyExit(err)
	}

	err = copyDataToFile(data, targetFile)
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
