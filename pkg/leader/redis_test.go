package leader_test

import (
	"log/slog"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func TestNewRedisLeader(t *testing.T) {
	t.Parallel()
	// start miniredis instance
	mrInstance := miniredis.RunT(t)

	// Create redis client
	mrClient := redis.NewClient(&redis.Options{
		Addr: mrInstance.Addr(),
	})

	ldr, err := leader.NewRedisLeader(mrClient, "leader:uuid")
	if err != nil {
		slog.Error("error creating UUID: %w", err)
	}

	slog.Debug("New Leader", "contains", ldr)
}
