package main

import (
	"fmt"
	"time"
)

func doAuth() error {
	// migrations
	dbType := micro.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().Unix())
	upFile := micro.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := micro.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyTemplateFile("templates/migrations/auth_tables."+dbType+".sql", upFile)
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

	err = copyTemplateFile("templates/data/user.go.txt", micro.RootPath+"/data/user.go")
	if err != nil {
		gracefullyExit(err)
	}

	err = copyTemplateFile("templates/data/token.go.txt", micro.RootPath+"/data/token.go")
	if err != nil {
		gracefullyExit(err)
	}

	return nil
}
