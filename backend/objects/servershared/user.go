package servershared

import "errors"

type User struct {
	Username     string
	Email        string
	Password     string
	PasswordHash []byte
	Salt         []byte
}

var ErrUsernameAlreadyExists = errors.New("username already exists")
