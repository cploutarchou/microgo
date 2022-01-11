package main

import (
	"github.com/joho/godotenv"
	"os"
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
