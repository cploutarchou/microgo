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
	Users  User
	Tokens Token
}

func New(dbPool *sql.DB) Models {
	db = dbPool
	switch os.Getenv("DATABASE_TYPE") {
	case "mysql", "mariadb":
		upper, _ = mysql.New(db)
	case "postgres", "postgresql":
		upper, _ = postgresql.New(db)
	default:
		// do Nothing

	}

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertID(i _db.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
