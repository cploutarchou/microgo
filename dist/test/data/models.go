package data

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

var db *bun.DB
var upper _db.Session

type Models struct {
	// Any model inserted here
	// Can be accessed through entire testlication

}

func New(dbPool *sql.DB) Models {
	switch os.Getenv("DATABASE_TYPE") {
	case "mysql", "mariadb":
		db = bun.NewDB(dbPool, mysqldialect.New())
	case "postgres", "postgresql":
		db = bun.NewDB(dbPool, pgdialect.New())
	default:
		db = bun.NewDB(dbPool, sqlitedialect.New())

	}

	return Models{}
}

// func getInsertID(i _db.ID) int {
// 	idType := fmt.Sprintf("%T", i)
// 	if idType == "int64" {
// 		return int(i.(int64))
// 	}

// 	return i.(int)
// }
