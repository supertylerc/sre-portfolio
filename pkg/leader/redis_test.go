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
	// start 3 miniredis instances
	mr1 := miniredis.RunT(t)
	mr2 := miniredis.RunT(t)
	mr3 := miniredis.RunT(t)

	// Create redis client
	mrClient1 := redis.NewClient(&redis.Options{
		Addr:       mr1.Addr(),
		ClientName: "Client 1",
	})

	mrClient2 := redis.NewClient(&redis.Options{
		Addr:       mr2.Addr(),
		ClientName: "Client 2",
	})

	mrClient3 := redis.NewClient(&redis.Options{
		Addr:       mr3.Addr(),
		ClientName: "Client 3",
	})

	ldr1, err := leader.NewRedisLeader(mrClient1, "leader:uuid")
	if err != nil {
		slog.Error("error creating UUID: %w", err)
	}

	ldr2, err := leader.NewRedisLeader(mrClient2, "leader:uuid")
	if err != nil {
		slog.Error("Error creating UUID: %w", err)
	}

	ldr3, err := leader.NewRedisLeader(mrClient3, "leader:uuid")
	if err != nil {
		slog.Error("Error creating UUID: %w", err)
	}

	slog.Info("New Leader", "contains", ldr1)
	slog.Info("New Leader", "contains", ldr2)
	slog.Info("New Leader", "contains", ldr3)

	ldr1.WriteLeader()
	ldr2.WriteLeader()
	ldr3.WriteLeader()

}
