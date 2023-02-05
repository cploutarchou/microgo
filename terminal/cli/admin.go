package main

import (
	"fmt"
)

func createAdminUser() error {
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

	dbType := micro.DB.DatabaseType

	if dbType == "mariadb" {
		dbType = "mysql"
	}

	if dbType == "postgresql" {
		dbType = "postgres"
	}

	// create the required tables in the database for roles , permissions and users if they don't exist

	_, err = micro.DB.Pool.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL,
		email VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	);`)
	if err != nil {
		return err
	}

	_, err = micro.DB.Pool.Exec(`CREATE TABLE IF NOT EXISTS roles (
		id SERIAL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP NULL,
		CONSTRAINT roles_pkey PRIMARY KEY (id)
	);`)
	if err != nil {
		return err
	}

	_, err = micro.DB.Pool.Exec(`CREATE TABLE IF NOT EXISTS permissions (
		id SERIAL,
		name VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP NULL,
		CONSTRAINT permissions_pkey PRIMARY KEY (id)
	);`)
	if err != nil {
		return err
	}

	_, err = micro.DB.Pool.Exec(`CREATE TABLE IF NOT EXISTS role_permissions (
		id SERIAL,
		role_id INT NOT NULL,
		permission_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP NULL,
		CONSTRAINT role_permissions_pkey PRIMARY KEY (id)
	);`)
	if err != nil {
		return err
	}

	_, err = micro.DB.Pool.Exec(`CREATE TABLE IF NOT EXISTS user_roles (
		id SERIAL,
		user_id INT NOT NULL,
		role_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP NULL,
		CONSTRAINT user_roles_pkey PRIMARY KEY (id)
	);`)
	if err != nil {
		return err
	}

	// insert the default roles if they don't exist
	_, err = micro.DB.Pool.Exec(`INSERT INTO roles (name, created_at, updated_at) VALUES ('admin', NOW(), NOW()) ON CONFLICT (name) DO NOTHING;`)
	if err != nil {
		return err
	}

	

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
