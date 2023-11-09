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
	miniInstance := miniredis.RunT(t)

	// Create new redis client context
	mrClient := redis.NewClient(&redis.Options{
		Addr:       miniInstance.Addr(),
		ClientName: "client1",
	})

	ldr, err := leader.NewRedisLeader(mrClient, "leader:uuid")
	if err != nil {
		slog.Error("error creating UUID: %w", err)
	}

	slog.Info("New Leader", "contains", ldr)

	// methods
	err = ldr.WriteLeader()

	if err != nil {
		slog.Error("WriteLeader() failed", "message", err)
	}

	isCur, err := ldr.IsCurrentLeader()
	if err != nil {
		slog.Info("IsCurrentLeader() failed", "message", err)
	}

	slog.Info("IsCurrentLeader()", "Result is", isCur)

	readRes, err := ldr.ReadLeader()
	if err != nil {
		slog.Info("ReadLeader() failed", "message", err)
	}

	slog.Info("ReadLeader()", "UUID", readRes)

	miniInstance.Close()
}
