package data

import (
	up "github.com/upper/db/v4"
	"time"
)

type User struct {
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
}

func (u *User) Table() string {
	return "users"
}

func (u *User) GetAll(condition up.Cond) ([]*User, error) {
	collection := upper.Collection(u.Table())
	var all []*User
	res := collection.Find(condition)
	err := res.All(&all)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (u *User) GetByID(id int) (*User, error) {
	var _user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{
		"id =": id,
	})
	err := res.One(&_user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())
	res = collection.Find(up.Cond{"user_id =": _user.ID, "expiry <": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}
	_user.Token = token
	return &_user, nil
}
