package main

import (
	"errors"
	"fmt"
	"time"
)

func makeDo(arg2, arg3 string) error {

	switch arg2 {
	case "migration":
		databaseType := micro.DB.DatabaseType
		if arg3 == "" {
			gracefullyExit(errors.New("something went wrong. Migration name is not specified"))
		}

		filename := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)
		migrationsUpFile := micro.RootPath + "/migrations/" + filename + "." + databaseType + ".up.sql"
		migrationsDownFile := micro.RootPath + "/migrations/" + filename + "." + databaseType + ".down.sql"
		err := copyTemplateFile("templates/migrations/migration."+databaseType+".up.sql", migrationsUpFile)
		if err != nil {
			gracefullyExit(err)
		}
		err = copyTemplateFile("templates/migrations/migration."+databaseType+".down.sql", migrationsDownFile)
		if err != nil {
			gracefullyExit(err)
		}
	}
	return nil
}
