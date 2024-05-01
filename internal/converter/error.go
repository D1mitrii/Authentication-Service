package converter

import "errors"

var (
	ErrEmptyEmail    = errors.New("empty email provided")
	ErrEmptyPassword = errors.New("empty password provided")
)
