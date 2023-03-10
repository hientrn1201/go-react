package models

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	Lastname  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdateAt  time.Time `json:"-"`
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		// check if the error is wrong password and not the server error
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
