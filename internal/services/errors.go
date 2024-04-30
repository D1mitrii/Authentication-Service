package services

import "errors"

var (
	ErrUserAlreadyExist  = errors.New("user already exist")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("incorrect user password")

	ErrCannotSignToken = errors.New("cannot sign token")
)
