//go:build integration

// run tests with this command: go test . --tags integration --count=1

package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secretpassword"
	dbName   = "microGo_test"
	port     = "5439"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=Europe/Athens connect_timeout=5"
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
var resource *dockertest.Resource
var pool *dockertest.Pool

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Something went wrong, unable not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Something went wrong, could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("Something went wrong, could not connect to docker: %s", err)
	}

	err = createTables(testDB)
	if err != nil {
		log.Fatalf("Something went wrong. Unable to create tables: %s", err)
	}

	models = New(testDB)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Something went wrong , could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables(db *sql.DB) error {
	stmt := `
	CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop table if exists users cascade;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    user_active integer NOT NULL DEFAULT 0,
    email character varying(255) NOT NULL UNIQUE,
    password character varying(60) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists remember_tokens;

CREATE TABLE remember_tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    remember_token character varying(100) NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON remember_tokens
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists tokens;

CREATE TABLE tokens (
    id SERIAL PRIMARY KEY,
    user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    first_name character varying(255) NOT NULL,
    email character varying(255) NOT NULL,
    token character varying(255) NOT NULL,
    token_hash bytea NOT NULL,
    created_at timestamp without time zone NOT NULL DEFAULT now(),
    updated_at timestamp without time zone NOT NULL DEFAULT now(),
    expiry timestamp without time zone NOT NULL
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON tokens
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
	`

	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}
	return nil
}

func TestUser_Table(t *testing.T) {
	s := models.Users.Table()
	if s != "users" {
		t.Error("Something went wrong, No valid table name returned: ", s)
	}
}

func TestUser_Insert(t *testing.T) {
	id, err := models.Users.Insert(dummyUSER)
	if err != nil {
		t.Error("Something went wrong, failed to insert user: ", err)
	}

	if id == 0 {
		t.Error("Something went wrong,  0 returned as id after insert")
	}
}

func TestUser_Get(t *testing.T) {
	u, err := models.Users.GetByID(1)
	if err != nil {
		t.Error("Something went wrong, unable to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("Something went wrong, returned user id  0 : ", err)
	}
}

func TestUser_GetAll(t *testing.T) {
	_, err := models.Users.GetAll()
	if err != nil {
		t.Error("Something went wrong, unable to get user: ", err)
	}
}

func TestUser_GetByEmail(t *testing.T) {
	u, err := models.Users.GetByEmail("cploutarchou@gmail.com")
	if err != nil {
		t.Error("Something went wrong, unable to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("Something went wrong, returned user id  0: ", err)
	}
}

func TestUser_Update(t *testing.T) {
	u, err := models.Users.GetByID(1)
	if err != nil {
		t.Error("Something went wrong, unable to get user: ", err)
	}

	u.LastName = "Kiriakou"
	err = u.Update(*u)
	if err != nil {
		t.Error("Something went wrong, unable to update user: ", err)
	}

	u, err = models.Users.GetByID(1)
	if err != nil {
		t.Error("Something went wrong, unable to get user: ", err)
	}

	if u.LastName != "Kiriakou" {
		t.Error("Something went wrong, unable to update user last name.")
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	u, err := models.Users.GetByID(1)
	if err != nil {
		t.Error("Something went wrong, unable to get user: ", err)
	}

	valid, err := u.VerifyPassword("mypassword")
	if err != nil {
		t.Error("Something went wrong, unable to verify user password: ", err)
	}

	if !valid {
		t.Error("Something went wrong, NOT valid user password.")
	}

	valid, err = u.VerifyPassword("mypassword")
	if err != nil {
		t.Error("Something went wrong, unable to verify user password: ", err)
	}

	if !valid {
		t.Error("Something went wrong, NOT valid user password.")
	}
}

func TestUser_ResetPassword(t *testing.T) {
	err := models.Users.ResetPassword(1, "new_pass")
	if err != nil {
		t.Error("Something went wrong, Unable to reset user password: ", err)
	}
	err = models.Users.ResetPassword(2, "new_pass")
	if err == nil {
		t.Error("Something went wrong, No any error while resetting user password for no valid user id. ", err)
	}

}

func TestUser_Delete(t *testing.T) {
	err := models.Users.Delete(1)
	if err != nil {
		t.Error("Something went wrong, unable to delete user: ", err)
	}

	_, err = models.Users.GetByID(1)
	if err == nil {
		t.Error("Something went wrong. Retrieved user that was actually deleted. ")
	}
}
func TestToken_Table(t *testing.T) {
	s := models.Tokens.Table()
	if s != "tokens" {
		t.Error("Something went wrong, unexpected table name returned for tokens")
	}
}

func TestToken_GenerateToken(t *testing.T) {
	id, err := models.Users.Insert(dummyUSER)
	if err != nil {
		t.Error("Something went wrong, Unable to create user: ", err)
	}

	_, err = models.Tokens.GenerateToken(id, time.Hour*24*365)
	if err != nil {
		t.Error("Something went wrong, error generating a new token: ", err)
	}
}

func TestToken_Insert(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUSER.Email)
	if err != nil {
		t.Error("Something went wrong, Unable to get user")
	}

	token, err := models.Tokens.GenerateToken(u.ID, time.Hour*24*365)
	if err != nil {
		t.Error("Something went wrong, Unable to generate token: ", err)
	}

	err = models.Tokens.Insert(*token, *u)
	if err != nil {
		t.Error("Something went wrong, Unable to insert token : ", err)
	}
}

func TestToken_GetUserForToken(t *testing.T) {
	token := "abc"
	_, err := models.Tokens.GetUserByToken(token)
	if err == nil {
		t.Error("Something went wrong. Expected an error but not received when getting user with a not valid token")
	}

	u, err := models.Users.GetByEmail(dummyUSER.Email)
	if err != nil {
		t.Error("failed to get user")
	}

	_, err = models.Tokens.GetByToken(u.Token.Text)
	if err != nil {
		t.Error("Something went wrong, Unable to get user with valid token: ", err)
	}
}

func TestToken_GetTokensForUser(t *testing.T) {
	tokens, err := models.Tokens.GetUserToken(1)
	if err != nil {
		t.Error(err)
	}

	if len(tokens) > 0 {
		t.Error("Something went wrong, tokens returned for non-existent user")
	}
}

func TestToken_Get(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUSER.Email)
	if err != nil {
		t.Error("Something went wrong, unable to get user")
	}

	_, err = models.Tokens.Get(u.Token.ID)
	if err != nil {
		t.Error("Something went wrong, unable to get user by token  id: ", err)
	}
}

func TestToken_GetByToken(t *testing.T) {
	u, err := models.Users.GetByEmail(dummyUSER.Email)
	if err != nil {
		t.Error("Something went wrong, unable to get user")
	}

	_, err = models.Tokens.GetByToken(u.Token.Text)
	if err != nil {
		t.Error("Something went wrong, Unable to get token by token: ", err)
	}

	_, err = models.Tokens.GetByToken("123")
	if err == nil {
		t.Error("Something went wrong, no error getting non-existing token by token: ", err)
	}
}
