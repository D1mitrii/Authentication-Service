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

type Hasher interface {
	Hash(string) (string, error)
	Compare(string, string) bool
}

type Services struct {
	JWT    JWT
	hasher Hasher
	repo   *repository.Repositories
}

func New(jwt JWT, hasher Hasher, repo *repository.Repositories) *Services {
	return &Services{
		JWT:    jwt,
		hasher: hasher,
		repo:   repo,
	}
}

func (s *Services) Register(ctx context.Context, user models.User) (int, error) {
	// TODO: Add user info validation

	hash, err := s.hasher.Hash(user.Password)
	if err != nil {
		return 0, ErrHashing
	}
	user.Password = hash
	return s.repo.User.CreateUser(ctx, user)
}

func (s *Services) Login(ctx context.Context, user models.User) (models.Token, error) {
	userFromDB, err := s.repo.User.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return models.Token{}, ErrUserNotFound
		}
		return models.Token{}, err
	}

	if !s.hasher.Compare(user.Password, userFromDB.Password) {
		return models.Token{}, ErrIncorrectPassword
	}

	return s.generateJWT(ctx, userFromDB)
}

func (s *Services) Logout(ctx context.Context, refreshToken string) error {
	_, err := s.repo.RefreshSession.GetSession(ctx, refreshToken)
	if err != nil {
		return err
	}
	go s.repo.RefreshSession.DeleteSession(ctx, refreshToken)
	return nil
}

func (s *Services) generateJWT(ctx context.Context, user models.User) (models.Token, error) {
	access, errAccess := s.JWT.NewAccessToken(user)
	refresh, errRefresh := s.JWT.NewRefreshToken()
	if errAccess != nil || errRefresh != nil {
		return models.Token{}, ErrCannotSignToken
	}
	go s.repo.RefreshSession.CreateSession(ctx, refresh, user.Id)
	return models.Token{Access: access, Refresh: refresh}, nil
}

func (s *Services) RefreshSession(ctx context.Context, refreshToken string) (models.Token, error) {
	userId, err := s.repo.RefreshSession.GetSession(ctx, refreshToken)
	if err != nil {
		return models.Token{}, nil
	}
	go s.repo.RefreshSession.DeleteSession(ctx, refreshToken)

	user, err := s.repo.User.GetUserById(ctx, userId)
	if err != nil {
		return models.Token{}, nil
	}
	return s.generateJWT(ctx, user)
}
