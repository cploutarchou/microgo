package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	up "github.com/upper/db/v4"
	"net/http"
	"strings"
	"time"
)

type Token struct {
	ID        int       `db:"id,omitempty"`
	UserID    int       `db:"user_id" json:"user_id"`
	FirstName string    `db:"first_name" json:"first_name"`
	Email     string    `db:"email" json:"email"`
	Text      string    `db:"token" json:"text"`
	Hash      []byte    `db:"token_hash" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	Expires   time.Time `db:"expiry" json:"expiry"`
}

func (t *Token) Table() string {
	return "tokens"
}

func (t *Token) Get(id int) (*Token, error) {

	var token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id": id})
	err := res.One(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
func (t *Token) GetByToken(token string) (*Token, error) {

	var _token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": token})
	err := res.One(&_token)
	if err != nil {
		return nil, err
	}

	return &_token, nil
}

func (t *Token) GetUserByToken(token string) (*User, error) {
	var u User
	var _token Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": token})
	err := res.One(&_token)
	if err != nil {
		return nil, err
	}
	collection = upper.Collection("users")
	res = collection.Find(up.Cond{"id": _token.UserID})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}
	u.Token = _token
	return &u, nil
}

func (t *Token) GetUserToken(userID int) ([]*Token, error) {
	var tokens []*Token
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id": userID})
	err := res.All(&tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (t *Token) Insert(token Token, u User) error {
	collection := upper.Collection(t.Table())

	// delete existing tokens
	res := collection.Find(up.Cond{"user_id": u.ID})
	err := res.Delete()
	if err != nil {
		return err
	}

	token.CreatedAt = time.Now()
	token.UpdatedAt = time.Now()
	token.FirstName = u.FirstName
	token.Email = u.Email

	_, err = collection.Insert(token)
	if err != nil {
		return err
	}

	return nil
}
func (t *Token) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id": id})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) DeleteByToken(token string) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"token": token})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID:  userID,
		Expires: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Text = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Text))
	token.Hash = hash[:]

	return token, nil
}

func (t *Token) Authenticate(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("No authorization header received. ")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("No authorization header received. ")
	}

	token := headerParts[1]

	if len(token) != 26 {
		return nil, errors.New("Token wrong size. ")
	}

	_token, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("No matching token found. ")
	}

	if _token.Expires.Before(time.Now()) {
		return nil, errors.New("Expired token. ")
	}

	user, err := t.GetUserByToken(token)
	if err != nil {
		return nil, errors.New("No matching user found. ")
	}

	return user, nil
}

func (t *Token) IsValid(token string) (bool, error) {
	user, err := t.GetUserByToken(token)
	if err != nil {
		return false, errors.New("No matching user found. ")
	}

	if user.Token.Text == "" {
		return false, errors.New("No matching token found. ")
	}

	if user.Token.Expires.Before(time.Now()) {
		return false, errors.New("Expired token. ")
	}

	return true, nil
}
