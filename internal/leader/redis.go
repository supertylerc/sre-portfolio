package leader

import (
	"fmt"

	"github.com/sagikazarmark/slog-shim"
	"github.com/spf13/viper"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func CreateRedisLeader() (*leader.RedisLeader, error) {
	redisHost := viper.Get("REDIS_HOST").(string)
	redisPort := viper.Get("REDIS_PORT").(string)
	redisPassword := viper.Get("REDIS_PASSWORD").(string)
	redisLeaderKey := viper.Get("REDIS_LEADER_KEY").(string)
	slog.Debug(
		"Creating leader",
		slog.String("redisHost", redisHost),
		slog.String("redisPort", redisPort),
		slog.String("redisPassword", redisPassword),
		slog.String("redisLeaderKey", redisLeaderKey),
	)

	ldr, err := leader.NewRedisLeader(
		fmt.Sprintf("%s:%s", redisHost, redisPort),
		redisPassword,
		redisLeaderKey,
	)
	if err != nil {
		return &leader.RedisLeader{}, fmt.Errorf("error creating a Redis Leader: %w", err)
	}

	slog.Info(
		"Created new leader",
		slog.String("uuid", ldr.UUID.String()),
		slog.String("key", ldr.Key),
	)

	return ldr, nil
}
