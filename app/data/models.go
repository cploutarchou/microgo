package data

import (
	"database/sql"
	"fmt"
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
		upper, _ = mysql.New(db)
	} else {
		upper, _ = postgresql.New(db)
	}
	return Models{}
}

func getInsertID(i _db.ID) int {

	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}
	return i.(int)
}
