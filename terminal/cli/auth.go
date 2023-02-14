package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func doAuth() error {
	// migrations
	dbType := micro.DB.DatabaseType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
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
	//Copy over data/models
	err = copyTemplateFile("templates/data/user.go.txt", micro.RootPath+"/data/user.go")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/data/token.go.txt", micro.RootPath+"/data/token.go")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/data/remember_token.go.txt", micro.RootPath+"/data/remember_token.go")
	if err != nil {
		gracefullyExit(err)
	}

	//Copy over middleware
	err = copyTemplateFile("templates/middleware/auth.go.txt", micro.RootPath+"/middleware/auth.go")
	if err != nil {
		gracefullyExit(err)
	}

	err = copyTemplateFile("templates/middleware/auth-token.go.txt", micro.RootPath+"/middleware/auth-token.go")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/middleware/remember.go.txt", micro.RootPath+"/middleware/remember.go")
	if err != nil {
		gracefullyExit(err)
	}
	//Copy over handlers
	err = copyTemplateFile("templates/handlers/auth-handlers.go.txt", micro.RootPath+"/handlers/auth-handlers.go")
	if err != nil {
		gracefullyExit(err)
	}

	//Copy over the views

	err = copyTemplateFile("templates/mailer/reset-password.html.tmpl", micro.RootPath+"/mail/reset-password.html.tmpl")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/mailer/reset-password.plain.tmpl", micro.RootPath+"/mail/reset-password.plain.tmpl")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/views/login.html", micro.RootPath+"/views/login.html")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/views/forgot.html", micro.RootPath+"/views/forgot.html")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/views/reset-password.html", micro.RootPath+"/views/reset-password.html")
	if err != nil {
		gracefullyExit(err)
	}

	color.Yellow(" - users, tokens, and remember_tokens migrations successfully created and executed")
	color.Yellow(" - user and tokens models successfully created")
	color.Yellow(" - auth middleware successfully created")
	color.Yellow("")
	color.Red("Don't forget to add user and tokens models in data/models.go, and add the appropriate " +
		"middleware to your Routes!")
	return nil
}
