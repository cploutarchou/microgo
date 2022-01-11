package microGo

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func (m *MicroGo) MigrateUp(dsn string) error {
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mi *migrate.Migrate) {
		_, _ = mi.Close()
	}(mig)

	if err = mig.Up(); err != nil {
		log.Println("Something went wrong during migration. Error:", err)
		return err
	}
	return nil
}

func (m *MicroGo) MigrateDownAll(dsn string) error {
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mi *migrate.Migrate) {
		_, _ = mi.Close()
	}(mig)

	if err := mig.Down(); err != nil {
		return err
	}

	return nil
}

func (m *MicroGo) Steps(n int, dsn string) error {
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mi *migrate.Migrate) {
		_, _ = mi.Close()
	}(mig)

	if err := mig.Steps(n); err != nil {
		return err
	}

	return nil
}

func (m *MicroGo) ForceMigrate(dsn string) error {
	mig, err := migrate.New("file://"+m.RootPath+"/migrations", dsn)
	if err != nil {
		return err
	}
	defer func(mi *migrate.Migrate) {
		_, _ = mi.Close()
	}(mig)

	if err := mig.Force(-1); err != nil {
		return err
	}

	return nil
}
