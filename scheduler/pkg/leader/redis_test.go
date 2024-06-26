package leader_test

import (
	"log/slog"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/supertylerc/sre-portfolio/scheduler/pkg/leader"
)

func TestNewRedisLeader(t *testing.T) {
	t.Parallel()

	ctx := leader.Ctx

	// Run miniredis server
	miniServer := miniredis.RunT(t)
	defer miniServer.Close()

	// Set up miniredis client no leader
	miniClient := redis.NewClient(&redis.Options{
		Addr:       miniServer.Addr(),
		ClientName: "testclient1",
	})

	// Set key with 100ms TTL
	err := miniClient.Set(ctx, "leader:uuid", "key", leader.LeaderTTL).Err()
	if err != nil {
		t.Errorf("Failed to set initial key/value")
	}

	miniServer.FastForward(leader.LeaderTTL)

	// Set up leader client on miniClient with RedisLeader Options{}

	ldr, err := leader.NewRedisLeader(miniClient, "leader:uuid")
	if err != nil {
		t.Errorf("NewRedisLeader() failed %v", err)
	}

	slog.Info("New Leader", "contains", ldr.UUID)

	err = ldr.WriteLeader()
	if err != nil {
		t.Errorf("WriteLeader() failed %v", err)
	}

	isLdr, err := ldr.IsCurrentLeader()
	if err != nil || isLdr == false {
		t.Errorf("Call to IsCurrentLeader() failed %v", err)
	}

	slog.Info("Current Leader", "status", isLdr)
	slog.Info("Current Leader", "UUID", ldr.UUID)
}
