package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func setup() {
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
	return "mysql://" + micro.BuildDataSourceName()
}

func help() {
	color.Yellow(`Available commands:

	help                  - show the help commands
	version               - print application version
	migrate               - runs all up migrations that have not been run previously
	migrate down          - reverses the most recent migration
	migrate reset         - runs all down migrations in reverse order, and then all up migrations
	make migration <name> - creates two new up and down migrations in the migrations folder
	
	`)
}
