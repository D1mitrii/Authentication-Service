package services

import (
	"auth/internal/models"
	"auth/internal/repository"
	"auth/internal/repository/repoerrors"
	"context"
	"errors"
	"log/slog"
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
	log    *slog.Logger
	JWT    JWT
	hasher Hasher
	repo   *repository.Repositories
}

func New(log *slog.Logger, jwt JWT, hasher Hasher, repo *repository.Repositories) *Services {
	return &Services{
		log:    log,
		JWT:    jwt,
		hasher: hasher,
		repo:   repo,
	}
}

func (s *Services) Register(ctx context.Context, user models.User) (int, error) {
	// TODO: Add user info validation
	const op = "Services.Register"
	log := s.log.With(
		slog.String("operation", op),
		slog.String("email", user.Email),
	)
	hash, err := s.hasher.Hash(user.Password)
	if err != nil {
		log.Info("invalid password")
		return 0, ErrHashing
	}
	user.Password = hash
	id, err := s.repo.User.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repoerrors.ErrAlreadyExist) {
			log.Info(err.Error())
			return 0, ErrUserAlreadyExist
		}
		log.Warn("failed to create user", slog.String("error", err.Error()))
		return 0, err
	}
	return id, nil
}

func (s *Services) Login(ctx context.Context, user models.User) (models.Token, error) {
	const op = "Services.Login"
	log := s.log.With(
		slog.String("operation", op),
		slog.String("email", user.Email),
	)
	userFromDB, err := s.repo.User.GetUserByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			log.Warn("user not found", slog.String("error", err.Error()))
			return models.Token{}, ErrUserNotFound
		}
		log.Error("failed to get user", slog.String("error", err.Error()))
		return models.Token{}, err
	}

	if !s.hasher.Compare(user.Password, userFromDB.Password) {
		log.Info("invalid password")
		return models.Token{}, ErrIncorrectPassword
	}

	return s.generateJWT(ctx, userFromDB)
}

func (s *Services) Logout(ctx context.Context, refreshToken string) error {
	const op = "Services.Logout"
	log := slog.With(
		slog.String("operation", op),
		slog.With("refresh-token", refreshToken),
	)
	if _, err := s.repo.RefreshSession.DeleteSession(ctx, refreshToken); err != nil {
		log.Error("failed to get-delete refresh session", slog.String("error", err.Error()))
		return err
	}
	log.Info("success logout")
	return nil
}

func (s *Services) generateJWT(ctx context.Context, user models.User) (models.Token, error) {
	const op = "Services.generateJWT"
	log := s.log.With(
		slog.String("operation", op),
		slog.String("email", user.Email),
	)
	access, errAccess := s.JWT.NewAccessToken(user)
	refresh, errRefresh := s.JWT.NewRefreshToken()
	if errAccess != nil || errRefresh != nil {
		log.Warn("failed to sign token")
		return models.Token{}, ErrCannotSignToken
	}

	if err := s.repo.RefreshSession.CreateSession(ctx, refresh, user.Id); err != nil {
		return models.Token{}, ErrSessionCreateFail
	}
	return models.Token{Access: access, Refresh: refresh}, nil
}

func (s *Services) RefreshSession(ctx context.Context, refreshToken string) (models.Token, error) {
	const op = "Services.RefreshSession"
	log := s.log.With(
		slog.String("operation", op),
		slog.String("refresh-token", refreshToken),
	)
	userId, err := s.repo.RefreshSession.DeleteSession(ctx, refreshToken)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			log.Warn(err.Error())
			return models.Token{}, ErrSessionNotFound
		}
		log.Error("failed to get-delete refresh session", slog.String("error", err.Error()))
		return models.Token{}, err
	}

	user, err := s.repo.User.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			log.Warn(err.Error())
			return models.Token{}, ErrUserNotFound
		}
		log.Warn("failed to get user", slog.String("error", err.Error()))
		return models.Token{}, nil
	}
	return s.generateJWT(ctx, user)
}
