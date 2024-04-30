package rdb

import (
	"auth/internal/repository/repoerrors"
	"context"
	"time"

	"github.com/go-redis/redis"
)

type RefreshSession struct {
	client      *redis.Client
	refresh_ttl time.Duration
}

func NewRefreshRepo(client *redis.Client, refresh_ttl time.Duration) *RefreshSession {
	return &RefreshSession{
		client:      client,
		refresh_ttl: refresh_ttl,
	}
}

func (r *RefreshSession) CreateSession(ctx context.Context, refreshToken string, id int) error {
	return r.client.Set(refreshToken, id, r.refresh_ttl).Err()
}

func (r *RefreshSession) GetSession(ctx context.Context, refreshToken string) (int, error) {
	id, err := r.client.Get(refreshToken).Int()
	if err == redis.Nil {
		return 0, repoerrors.ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *RefreshSession) DeleteSession(ctx context.Context, refreshToken string) error {
	if r.client.Del(refreshToken).Val() == 0 {
		return repoerrors.ErrNotFound
	}
	return nil
}
