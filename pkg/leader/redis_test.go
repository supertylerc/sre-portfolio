package leader_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func TestNewRedisLeader(t *testing.T) {
	t.Parallel()

	ctx := leader.Ctx

	// Run miniredis server
	miniInstance := miniredis.RunT(t)

	// Set up miniredis client no leader
	mrClient := redis.NewClient(&redis.Options{
		Addr:       miniInstance.Addr(),
		ClientName: "client1",
	})

	// Set key with 100ms TTL

	err := mrClient.Set(ctx, "key", "value", 100*time.Millisecond).Err()
	if err != nil {
		panic(err)
	}

	time.Sleep(leader.LeaderTTL)

	// Set up leader client on mrClient
	ldr, err := leader.NewRedisLeader(mrClient, "leader:uuid")
	if err != nil {
		slog.Error("error creating UUID: %w", err)
	}

	slog.Info("New Leader", "contains", ldr)

	err = ldr.WriteLeader()
	if err != nil {
		slog.Error("WriteLeader() failed", err)
	}

	isLdr, err := ldr.IsCurrentLeader()
	if err != nil {
		slog.Error("Error checking current leader", err)
	}

	slog.Info("Current Leader", "contains", isLdr)

	miniInstance.Close()
}
