package data

import (
	up "github.com/upper/db/v4"
	"time"
)

type RememberToken struct {
	ID            int       `db:"id,omitempty"`
	UserID        int       `db:"user_id"`
	RememberToken string    `db:"remember_token"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

func (r *RememberToken) Table() string {
	return "remember_token"
}

func (r *RememberToken) InsertToken(userID int, token string) error {
	collection := upper.Collection(r.Table())
	rememberToken := RememberToken{
		UserID:        userID,
		RememberToken: token,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	_, err := collection.Insert(rememberToken)
	if err != nil {
		return err
	}
	return nil
}

func (r *RememberToken) Delete(rememberToken string) error {
	collection := upper.Collection(r.Table())
	res := collection.Find(up.Cond{"remember_token": rememberToken})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}
