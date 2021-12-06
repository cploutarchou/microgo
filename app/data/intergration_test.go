package data

import (
	"database/sql"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"log"
	"os"
	"testing"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secretpassword"
	dbName   = "microGo_test"
	port     = "5439"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5"
)
var dummyUSER = User{
	FirstName: "Christos",
	LastName:  "Ploutarchou",
	Email:     "cploutarchou@gmail.com",
	Active:    1,
	Password:  "mypassword",
}
var models Models
var testDB *sql.DB
var resource = *dockertest.Resource
var pool = *dockertest.Pool

func TestMain(m *testing.M) {
	_ = os.Setenv("DATABASE_TYPE", "postgres")
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal("Unable to connect with docker. ")
	}
	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-bullseye",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		}, ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Could not start docker resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Unable connect to docker: %s", err)
	}
}
