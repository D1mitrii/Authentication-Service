package services

import (
	"auth/internal/models"
)

type JWT interface {
	NewAccessToken(models.User) (string, error)
	NewRefreshToken() (string, error)
	Parse(string) (int, error)
}

type Services struct {
	JWT JWT
}

func New(jwt JWT) *Services {
	return &Services{
		JWT: jwt,
	}
}
