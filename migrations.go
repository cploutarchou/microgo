package MicroGO

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	migrate "github.com/rubenv/sql-migrate"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (m *MicroGo) MigrateUp(dsn string) error {
	var path string
	db, err := m.DB.Client.DB()
	if err != nil {
		log.Println("Error connecting to database")
		return err

	}
	path = "file://" + m.RootPath + "/migrations"
	if runtime.GOOS == "windows" {
		path = fmt.Sprint(strings.Replace(path, "/", "\\", -1))
	} else {
		path = "file://" + m.RootPath + "/migrations"
	}
	migrations := migrate.FileMigrationSource{
		Dir: path,
	}
	n, err := migrate.Exec(db, m.DB.DatabaseType, migrations, migrate.Up)

	if err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	log.Printf("Applied %d migrations!\n", n)
	return nil
}

func (m *MicroGo) MigrateDownAll(dsn string) error {
	var path string
	db, err := m.DB.Client.DB()
	if err != nil {
		log.Println("Error connecting to database")
		return err

	}
	path = "file://" + m.RootPath + "/migrations"
	if runtime.GOOS == "windows" {
		path = fmt.Sprint(strings.Replace(path, "/", "\\", -1))
	} else {
		path = "file://" + m.RootPath + "/migrations"
	}
	migrations := migrate.FileMigrationSource{
		Dir: path,
	}
	n, err := migrate.Exec(db, m.DB.DatabaseType, migrations, migrate.Down)

	if err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	log.Printf("Applied %d migrations!\n", n)
	return nil
}
