/*
 * MIT License
 *
 * Copyright (c) 2022 Christos Ploutarchou
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"github.com/fatih/color"
	"time"
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
	err = copyTemplateFile("templates/views/login.jet", micro.RootPath+"/views/login.jet")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/views/forgot.jet", micro.RootPath+"/views/forgot.jet")
	if err != nil {
		gracefullyExit(err)
	}
	err = copyTemplateFile("templates/views/reset-password.jet", micro.RootPath+"/views/reset-password.jet")
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
