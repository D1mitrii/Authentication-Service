package pgdb

import (
	"auth/internal/models"
	"auth/internal/repository/repoerrors"
	"auth/pkg/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (int, error) {
	const op = "UserRepo.CreateUser"
	sql := `INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id;`
	var id int
	err := r.Pool.QueryRow(ctx, sql, user.Email, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrors.ErrAlreadyExist
			}
		}
		return 0, fmt.Errorf("%s - r.Pool.QueryRow: %v", op, err)
	}
	return id, nil
}

func (r *UserRepo) GetUserById(ctx context.Context, id int) (models.User, error) {
	const op = "UserRepo.GetUserById"
	sql := `SELECT (id, email, password, created_at) FROM users WHERE id = $1;`
	var user models.User
	err := r.Pool.QueryRow(ctx, sql, id).Scan(&user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, repoerrors.ErrNotFound
		}
		return models.User{}, fmt.Errorf("%s - r.Pool.QueryRow: %v", op, err)
	}
	return user, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "UserRepo.GetUserByEmail"
	sql := `SELECT (id, email, password, created_at) FROM users WHERE email = $1;`
	var user models.User
	err := r.Pool.QueryRow(ctx, sql, email).Scan(&user)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, repoerrors.ErrNotFound
		}
		return models.User{}, fmt.Errorf("%s - r.Pool.QueryRow: %v", op, err)
	}
	return user, nil
}

func (r *UserRepo) DeleteUser(ctx context.Context, id int) error {
	const op = "UserRepo.GetUserById"
	sql := `DELETE FROM users WHERE id = $1;`
	_, err := r.Pool.Exec(ctx, sql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repoerrors.ErrNotFound
		}
		return fmt.Errorf("%s - r.Pool.QueryRow: %v", op, err)
	}
	return nil
}
