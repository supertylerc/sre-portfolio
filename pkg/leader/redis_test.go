package leader_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func TestNewRedisLeader(t *testing.T) {

	miniServer := miniredis.RunT(t)
	defer miniServer.Close()

	// ctx := leader.Ctx

	tests := []struct {
		name        string
		input       string
		want        string
		miniClient  *redis.Client
		miniOptions *redis.Options
		leader      *leader.RedisLeader
		ldrCheck    bool
	}{
		{
			name:  "Become the leader",
			input: "leader:uuid",
			want:  "Becoming the leader.",
			miniClient: redis.NewClient(&redis.Options{
				Addr:       miniServer.Addr(),
				ClientName: "testclient1",
			}),
		},
		{
			name:  "Become non-leader",
			input: "",
			want:  "Becoming the leader.",
			miniClient: redis.NewClient(&redis.Options{
				Addr:       miniServer.Addr(),
				ClientName: "testclient1",
			}),
		},
	}

	// miniServer.FastForward(leader.LeaderTTL)
	t.Parallel()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ldr, _ := leader.NewRedisLeader(tc.miniClient, tc.input)
			result, _ := ldr.IsCurrentLeader()

			if result != tc.ldrCheck {
				t.Errorf("NewRedisLeader() %v = %v, want %v", tc.name, ldr, tc.want)
			}
		})
	}
}
