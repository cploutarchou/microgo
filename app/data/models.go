package data

import (
	"database/sql"
	_db "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
	"github.com/upper/db/v4/adapter/postgresql"
	"os"
)

var db *sql.DB
var upper _db.Session

type Models struct {
	// Any model inserted here
	// Can be accessed through entire application

}

func New(dbPool *sql.DB) Models {
	db = dbPool
	if os.Getenv("DATABASE_TYPE") != "mysql" || os.Getenv("DATABASE_TYPE") != "mariadb" {
		upper, _ = mysql.New(dbPool)
	} else {
		upper, _ = postgresql.New(dbPool)
	}
	return Models{}
}
