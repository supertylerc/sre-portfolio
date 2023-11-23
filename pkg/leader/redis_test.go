package leader_test

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func TestNewRedisLeader(t *testing.T) {
	t.Parallel()

	miniServer := miniredis.RunT(t)
	t.Cleanup(func() { miniServer.Close() })

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

	for _, testcase := range tests {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			t.Parallel()
			ldr, _ := leader.NewRedisLeader(testcase.miniClient, testcase.input)
			result, _ := ldr.IsCurrentLeader()

			if result != testcase.ldrCheck {
				t.Errorf("NewRedisLeader() %v = %v, want %v", testcase.name, ldr, testcase.want)
			}
		})
	}
}
