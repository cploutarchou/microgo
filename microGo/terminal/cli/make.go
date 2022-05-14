package main

import (
	"errors"
	"fmt"
	"time"
)

func makeDo(arg2, arg3 string) error {

	switch arg2 {

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

	}

	return nil
}
