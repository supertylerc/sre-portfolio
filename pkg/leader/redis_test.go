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

	t.Parallel()
	//ctx := leader.Ctx

	tests := []struct {
		name        string
		input       string
		want        string
		miniClient  *redis.Client
		miniOptions *redis.Options
		leader      *leader.RedisLeader
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
	}

	// miniServer.FastForward(leader.LeaderTTL)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := leader.NewRedisLeader(tc.miniClient, tc.input)

			if got != tc.leader {
				t.Errorf("NewRedisLeader() = %v, want %v", got, tc.want)
			}
		})
	}
}
