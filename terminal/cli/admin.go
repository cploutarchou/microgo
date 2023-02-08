package main

import (
	"fmt"
	"time"
)

func createAdminUser() error {
	askForInput("")
	dbType := micro.DB.DatabaseType
	
	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := micro.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := micro.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyTemplateFile("templates/migrations/admin_tables."+dbType+".sql", upFile)
	if err != nil {
		gracefullyExit(err)
	}

	err = copyDataToFile([]byte("drop table if exists users cascade; drop table if exists tokens cascade; drop table if exists remember_tokens;"), downFile)
	if err != nil {
		gracefullyExit(err)
	}

	// run migrations
	err = doMigrate("up", "")
	if err != nil {
		gracefullyExit(err)
	}
	// ask for email
	email, err := askForInput("Enter admin email: ")
	if err != nil {
		return err
	}
	// ask for password
	password, err := askForInput("Enter admin password: ")
	if err != nil {
		return err
	}
	// ask for password again
	password2, err := askForInput("Enter admin password again: ")
	if err != nil {
		return err
	}
	// check if passwords match
	if password != password2 {
		return fmt.Errorf("passwords do not match")
	}
	// create user

	fmt.Println("Creating admin user...")
	_, err = micro.DB.Pool.Exec("INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4)", email, password, time.Now(), time.Now())
	if err != nil {
		return err
	}
	fmt.Println("Admin user created successfully!")
	return nil

}

func askForInput(question string) (string, error) {
	fmt.Print(question)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	return input, nil
}
