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

	ctx := leader.Ctx

	miniServer := miniredis.RunT(t)
	defer miniServer.Close()

	// Set up miniredis client no leader
	miniClient := redis.NewClient(&redis.Options{
		Addr:       miniServer.Addr(),
		ClientName: "testclient1",
	})

	err := miniClient.Set(ctx, "leader:uuid", "key", leader.LeaderTTL).Err()
	if err != nil {
		t.Errorf("Failed to set initial key/value")
	}

	miniServer.FastForward(leader.LeaderTTL)

	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "Become the leader",
			input: "leader:uuid",
			want:  true,
		},
		{
			name:  "Not the leader",
			input: "",
			want:  true,
		},
	}

	for _, testcase := range tests {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			miniServer.FastForward(leader.LeaderTTL)
			ldr, _ := leader.NewRedisLeader(miniClient, testcase.input)
			slog.Info("New Leader", "contains", ldr.UUID)
			_ = ldr.WriteLeader()
			isLdr, _ := ldr.IsCurrentLeader()
			slog.Info("Current Leader", "status", isLdr)
			slog.Info("Current Leader", "UUID", ldr.UUID)

			if isLdr != testcase.want {
				t.Errorf("got %t; want %t", isLdr, testcase.want)
			}
		})
	}
}
