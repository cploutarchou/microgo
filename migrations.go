package MicroGO

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (m *MicroGo) MigrateUp(dsn string) error {
	var path string
	path = "file://" + m.RootPath + "/migrations"
	if runtime.GOOS == "windows" {
		path = fmt.Sprintf(strings.Replace(path, "/", "\\", -1))
	} else {
		path = "file://" + m.RootPath + "/migrations"
	}
	mig, err := migrate.New(path, dsn)
	if err != nil {
		return err
	}
	defer func(mig *migrate.Migrate) {
		_, _ = mig.Close()
	}(mig)

	if err := mig.Up(); err != nil {
		log.Println("Error running migration:", err)
		return err
	}
	return nil
}

func (m *MicroGo) MigrateDownAll(dsn string) error {
	// TODO: Add windows support.
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mig *migrate.Migrate) {
		_, _ = mig.Close()
	}(mig)

	if err := mig.Down(); err != nil {
		return err
	}

	return nil
}

func (m *MicroGo) Steps(n int, dsn string) error {
	// TODO: Add windows support.
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mig *migrate.Migrate) {
		_, _ = mig.Close()
	}(mig)

	if err := mig.Steps(n); err != nil {
		return err
	}

	return nil
}

func (m *MicroGo) MigrateForce(dsn string) error {
	// TODO: Add windows support.
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mig *migrate.Migrate) {
		_, _ = mig.Close()
	}(mig)

	if err := mig.Force(-1); err != nil {
		return err
	}

	return nil
}
