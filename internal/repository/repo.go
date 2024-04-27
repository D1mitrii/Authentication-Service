package repository

import (
	"auth/internal/models"
	"context"
)

type UserRepo interface {
	CreateUser(context.Context, models.User) (int, error)
	GetUserById(context.Context, int) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	DeleteUser(context.Context, int) error
}

type Repositories struct {
	User UserRepo
}

func New(users UserRepo) *Repositories {
	return &Repositories{
		User: users,
	}
}
