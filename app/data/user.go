package data

import (
	"errors"
	up "github.com/upper/db/v4"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) GetAll() ([]*User, error) {
	collection := upper.Collection(u.Table())
	var all []*User
	res := collection.Find().OrderBy("created_at")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	var _user User
	collection := upper.Collection(u.Table())
	res := collection.Find(up.Cond{
		"email =": email,
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

func (u *User) Update(user User) error {
	user.UpdatedAt = time.Now()
	collection := upper.Collection(u.Table())
	res := collection.Find(user.ID)

	err := res.Update(&user)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete(id User) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Insert(user User) (int, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = string(hashed)

	collection := upper.Collection(u.Table())
	res, err := collection.Insert(user)
	if err != nil {
		return 0, err
	}
	id := getInsertID(res.ID())
	return id, nil
}

func (u *User) ResetPassword(id int, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 16)
	if err != nil {
		return err
	}
	user, err := u.GetByID(id)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	err = user.Update(*u)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) VerifyPassword(text string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(text))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
