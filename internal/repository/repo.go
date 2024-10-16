package repository

import (
	"context"
	"github.com/d1mitrii/authentication-service/internal/models"
)

type UserRepo interface {
	CreateUser(context.Context, models.User) (int, error)
	GetUserById(context.Context, int) (models.User, error)
	GetUserByEmail(context.Context, string) (models.User, error)
	DeleteUser(context.Context, int) error
}

type RefreshSessionRepo interface {
	CreateSession(context.Context, string, int) error
	GetSession(context.Context, string) (int, error)
	DeleteSession(context.Context, string) (int, error)
}

type Repositories struct {
	User           UserRepo
	RefreshSession RefreshSessionRepo
}

func New(users UserRepo, session RefreshSessionRepo) *Repositories {
	return &Repositories{
		User:           users,
		RefreshSession: session,
	}
}
