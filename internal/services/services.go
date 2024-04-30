package services

import (
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/repository/repoerrors"
	"context"
	"errors"
)

type JWT interface {
	NewAccessToken(models.User) (string, error)
	NewRefreshToken() (string, error)
	Parse(string) (int, error)
}

type Services struct {
	JWT  JWT
	Repo *repository.Repositories
}

func New(jwt JWT, repo *repository.Repositories) *Services {
	return &Services{
		JWT:  jwt,
		Repo: repo,
	}
}

func (s *Services) Register(ctx context.Context, user models.User) (int, error) {
	// TODO: Add user info validation
	// TODO: Save user with hashed password
	return s.Repo.User.CreateUser(ctx, user)
}

func (s *Services) Login(ctx context.Context, user models.User) (models.JWTPair, error) {
	userFromDB, err := s.Repo.User.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return models.JWTPair{}, ErrUserNotFound
		}
		return models.JWTPair{}, err
	}
	// TODO: Hash req user password
	if userFromDB.Password != user.Password {
		return models.JWTPair{}, ErrIncorrectPassword
	}

	access, errAccess := s.JWT.NewAccessToken(userFromDB)
	refresh, errRefresh := s.JWT.NewRefreshToken()
	if errAccess != nil || errRefresh != nil {
		return models.JWTPair{}, ErrCannotSignToken
	}
	// TODO: Save refresh token to Redis: key - refresh, val - user id
	return models.JWTPair{Access: access, Refresh: refresh}, nil
}

func (s *Services) Logout(ctx context.Context, refreshToken string) error {
	return nil
}

func (s *Services) RefreshSession(ctx context.Context, refreshToken string) (models.JWTPair, error) {
	return models.JWTPair{}, nil
}
