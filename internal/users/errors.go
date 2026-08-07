package users

import "errors"

var (
	ErrNotFound      = errors.New("user not found")
	ErrInvalidInput  = errors.New("invalid user input")
	ErrUsernameTaken = errors.New("username already exists")
	ErrEmailTaken    = errors.New("email already exists")
)
