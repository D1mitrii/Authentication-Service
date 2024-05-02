package rdb

import (
	"auth/internal/repository/repoerrors"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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
	return r.client.Set(ctx, refreshToken, id, r.refresh_ttl).Err()
}

func (r *RefreshSession) GetSession(ctx context.Context, refreshToken string) (int, error) {
	const op = "RefreshSession.GetSession"
	id, err := r.client.Get(ctx, refreshToken).Int()
	if err == redis.Nil {
		return 0, repoerrors.ErrNotFound
	} else if err != nil {
		return 0, fmt.Errorf("%s - client.Get: %v", op, err)
	}
	return id, nil
}

func (r *RefreshSession) DeleteSession(ctx context.Context, refreshToken string) (int, error) {
	const op = "RefreshSession.DeleteSession"
	id, err := r.client.GetDel(ctx, refreshToken).Int()
	if err == redis.Nil {
		return 0, repoerrors.ErrNotFound
	} else if err != nil {
		return 0, fmt.Errorf("%s - client.Get: %v", op, err)
	}
	return id, nil
}
