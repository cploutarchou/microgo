package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

//go:embed templates
var templateFS embed.FS

func copyTemplateFile(templatePath, targetFile string) error {
	if fileExists(targetFile) {
		return errors.New(targetFile + " already exists!")
	}
	fileUrl := "https://raw.githubusercontent.com/cploutarchou/MicroGO/master/terminal/cli/"
	resp, err := http.Get(fileUrl + templatePath)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)

	if err != nil {
		gracefullyExit(err)
	}
	if strings.Contains(string(data), "$APPNAME$/data") {

		appName := strings.Split(micro.RootPath, "/")
		micro.AppName = strings.ReplaceAll(appName[len(appName)-1], "/data", "")
		err = copyDataToFile(
			[]byte(strings.ReplaceAll(string(data),
				"$APPNAME$/", fmt.Sprintf("%s/", micro.AppName))),
			targetFile,
		)
	} else {
		err = copyDataToFile(data, targetFile)
	}

	if err != nil {
		gracefullyExit(err)
	}

	return nil
}
func readFromRepo(fileToRead string) ([]byte, error) {
	fileUrl := "https://raw.githubusercontent.com/cploutarchou/MicroGO/master/terminal/cli/"
	resp, err := http.Get(fileUrl + fileToRead)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func copyDataToFile(data []byte, to string) error {
	err := os.WriteFile(to, data, 0644)
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
