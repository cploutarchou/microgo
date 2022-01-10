package data

import (
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	_db "github.com/upper/db/v4"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	fakeDB, _, _ := sqlmock.New()
	defer fakeDB.Close()

	_ = os.Setenv("DATABASE_TYPE", "postgres")
	m := New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("No valid type", fmt.Sprintf("%T", m))
	}

	_ = os.Setenv("DATABASE_TYPE", "mysql")
	m = New(fakeDB)
	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("No valid type", fmt.Sprintf("%T", m))
	}
}

func TestGetInsertID(t *testing.T) {
	var id _db.ID
	id = int64(1)

	returnedID := getInsertID(id)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("No valid type returned. ")
	}

	id = 1
	returnedID = getInsertID(id)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("No valid type returned. ")
	}

}
