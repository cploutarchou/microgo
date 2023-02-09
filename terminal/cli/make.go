package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func makeDo(arg2, arg3 string) error {

	switch arg2 {
	case "key":
		rnd := micro.CreateRandomString(32)
		color.Yellow("Successfully created a 32 chars encryption key :  %s", rnd)
	case "auth":
		err := doAuth()
		if err != nil {
			gracefullyExit(err)
		}
	case "migration":
		dbType := micro.DB.DatabaseType
		if arg3 == "" {
			gracefullyExit(errors.New("you must give the migration a name"))
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)

		upFile := micro.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		downFile := micro.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		err := copyTemplateFile("templates/migrations/migration."+dbType+".up.sql", upFile)
		if err != nil {
			gracefullyExit(err)
		}

		err = copyTemplateFile("templates/migrations/migration."+dbType+".down.sql", downFile)
		if err != nil {
			gracefullyExit(err)
		}
	case "handler":
		if arg3 == "" {
			gracefullyExit(errors.New("you must give the handler a name"))
		}

		fileName := micro.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			gracefullyExit(errors.New(fileName + " already exists!"))
		}

		data, err := readFromRepo("templates/handlers/handler.go.txt")
		if err != nil {
			gracefullyExit(err)
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLERNAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			gracefullyExit(err)
		}
	case "model":
		if arg3 == "" {
			gracefullyExit(errors.New("you must give a name to your model"))
		}
		data, err := readFromRepo("templates/data/model.go.txt")
		if err != nil {
			gracefullyExit(err)
		}
		model := string(data)
		prul := pluralize.NewClient()
		var modelName = arg3
		var tableName = arg3
		if prul.IsPlural(arg3) {
			modelName = prul.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(prul.Plural(arg3))
		}
		fileName := micro.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			gracefullyExit(errors.New(fileName + " already exists!"))
		}
		model = strings.ReplaceAll(model, "$MODELNAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLENAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			gracefullyExit(err)
		}
	case "session":
		err := createSessionTable()
		if err != nil {
			gracefullyExit(err)
		}
	case "mail":
		if arg3 == "" {
			gracefullyExit(errors.New("you must specify template file name! "))
		}
		htmlMail := micro.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainTextMail := micro.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"
		err := copyTemplateFile("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			gracefullyExit(err)
		}
		err = copyTemplateFile("templates/mailer/mail.plain.tmpl", plainTextMail)
		if err != nil {
			gracefullyExit(err)
		}
	case "admin":
		// create admin user
		err := createAdminUser()

		
		if err != nil {
			gracefullyExit(err)
		}
	}
	return nil
}
