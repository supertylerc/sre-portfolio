package leader

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// nolint: varnamelen
func CreateRedisClient(ctx context.Context, db int) (*redis.Client, error) {
	redisHost := viper.Get("REDIS_HOST").(string)
	redisPort := viper.Get("REDIS_PORT").(string)
	redisPassword := viper.Get("REDIS_PASSWORD").(string)

	slog.Debug(
		"Creating client",
		slog.String("redisHost", redisHost),
		slog.String("redisPort", redisPort),
		slog.String("redisPassword", redisPassword),
		slog.Int("redisDB", db),
	)

	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       db,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to ping Redis: %w", err)
	}

	slog.Info(
		"Created new redis client",
		slog.String("address", redisAddress),
		slog.Int("db", db),
	)

	return client, nil
}
