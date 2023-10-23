package leader

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// nolint: gochecknoglobals
var (
	ErrNil = errors.New("no matching record found")
	Ctx    = context.TODO()
)

const (
	LeaderTTL time.Duration = 1 * time.Second
)

type RedisLeader struct {
	Client *redis.Client
	UUID   uuid.UUID
	Key    string
}

func NewRedisLeader(client *redis.Client, key string) (*RedisLeader, error) {
	ldrUUID, err := uuid.NewUUID()
	if err != nil {
		return nil, fmt.Errorf("error creating UUID: %w", err)
	}

	return &RedisLeader{
		Client: client,
		UUID:   ldrUUID,
		Key:    key,
	}, nil
}

func (ldr *RedisLeader) WriteLeader() error {
	isCurrentLeader, err := ldr.IsCurrentLeader()

	switch {
	case err != nil && !errors.Is(err, ErrNil):
		return err
	case errors.Is(err, ErrNil):
		slog.Info(
			"Becoming the leader.",
			slog.String("uuid", ldr.UUID.String()),
			slog.String("key", ldr.Key),
		)

		if err = ldr.Client.Set(Ctx, ldr.Key, ldr.UUID.String(), LeaderTTL).Err(); err != nil {
			return fmt.Errorf("unable to set Redis leader: %w", err)
		}

		totalLeadershipClaimed.WithLabelValues(ldr.UUID.String()).Inc()
	case isCurrentLeader:
		slog.Debug(
			"Refreshing leader lock",
			slog.String("uuid", ldr.UUID.String()),
			slog.String("key", ldr.Key),
		)

		if err = ldr.Client.Expire(Ctx, ldr.Key, LeaderTTL).Err(); err != nil {
			return fmt.Errorf("unable to expire leader key: %w", err)
		}

		totalLeadershipRefreshed.WithLabelValues(ldr.UUID.String()).Inc()
	}

	return nil
}

func (ldr *RedisLeader) ReadLeader() (string, error) {
	res, err := ldr.Client.Get(Ctx, ldr.Key).Result()
	totalLeadershipRead.WithLabelValues(ldr.UUID.String()).Inc()

	if err != nil {
		return "", ErrNil
	}

	return res, nil
}

func (ldr *RedisLeader) IsCurrentLeader() (bool, error) {
	res, err := ldr.ReadLeader()
	if err != nil {
		return false, err
	}

	val := 0.0
	if res == ldr.UUID.String() {
		val = 1.0
	}

	currentLeader.WithLabelValues(ldr.UUID.String()).Set(val)

	return val == 1.0, nil
}
