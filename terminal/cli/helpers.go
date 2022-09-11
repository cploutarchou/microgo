package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup(arg1, arg2 string) {
	if arg1 != "new" && arg1 != "version" && arg1 != "help" {
		err := godotenv.Load()
		if err != nil {
			gracefullyExit(err)
		}

		path, err := os.Getwd()
		if err != nil {
			gracefullyExit(err)
		}

		micro.RootPath = path
		micro.DB.DatabaseType = os.Getenv("DATABASE_TYPE")
	}
}

func getDSN() string {
	dbType := micro.DB.DatabaseType

	if dbType == "pgx" {
		dbType = "postgres"
	}

	if dbType == "postgres" {
		var dsn string
		if os.Getenv("DATABASE_PASS") != "" {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_PASS"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		} else {
			dsn = fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
				os.Getenv("DATABASE_USER"),
				os.Getenv("DATABASE_HOST"),
				os.Getenv("DATABASE_PORT"),
				os.Getenv("DATABASE_NAME"),
				os.Getenv("DATABASE_SSL_MODE"))
		}
		return dsn
	}
	return "mysql://" + micro.BuildDSN()
}

func help() {
	color.Yellow(`Available commands:

	help                  - show the help commands
	version               - print application version
	make auth             - Create and runs migrations for auth tables, create models and middleware.
	migrate               - runs all up migrations that have not been run previously
	migrate down          - reverses the most recent migration
	migrate reset         - runs all down migrations in reverse order, and then all up migrations
	make migration <name> - creates two new up and down migrations in the migrations folder
	make handler <name>   - create a stub handler on handlers directory
	make model <name>     - create a new mode in the data directory
	make session          - create a new table in the database as a session storage
	make key              - create a random key of 32 characters.
	make mail             - create two starter mail templates in the mail directory.
	
	`)
}

func updateSrcFiles(path string, fi os.FileInfo, err error) error {
	// check for errors
	if err != nil {
		return err
	}
	// check if is dir
	if fi.IsDir() {
		return nil
	}
	// check if is go file
	match, err := filepath.Match("*.go", fi.Name())
	if err != nil {
		return err
	}
	if match {
		//read file
		read, err := os.ReadFile(path)
		if err != nil {
			gracefullyExit(err)
		}
		newContent := strings.Replace(string(read), "app", appURL, -1)
		// save file
		err = os.WriteFile(path, []byte(newContent), 0)
		if err != nil {
			gracefullyExit(err)
		}
	}
	return nil
}

// updateSrcFolders walks the given path and updates the src files and the sub folders
func updateSrcFolders() error {
	err := filepath.Walk(".", updateSrcFiles)
	if err != nil {
		return err
	}
	return nil
}
