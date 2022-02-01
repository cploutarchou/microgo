package main

import (
	"fmt"
	"log"
	"time"
)

func doAuth() error {

	//migrations

	dbType := micro.DB.DatabaseType
	filename := fmt.Sprintf("%d_create_auth_tables", time.Now().Unix())
	upFile := micro.RootPath + "/migrations" + filename + ".up.sql"
	downFile := micro.RootPath + "/migrations" + filename + ".down.sql"
	log.Println(dbType, upFile, downFile)
	err := copyTemplateFile("templates/migrations/auth_tables"+dbType+".sql", upFile)
	if err != nil {
		gracefullyExit(err)
	}
	err = copyDataToFile([]byte("drop table if exists users cascade"), downFile)
	if err != nil {
		gracefullyExit(err)
	}
	//run migrations
	err = doMigrate("up", "")
	if err != nil {
		gracefullyExit(err)
	}
	//copy files
	return nil
}
