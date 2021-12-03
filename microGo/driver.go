package microGo

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func (m *MicroGo) OpenDB(driverName, dataSourceName string) (*sql.DB, error) {
	if driverName == "postgres" || driverName == "postgresql" {
		driverName = "pgx"
	}

	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
