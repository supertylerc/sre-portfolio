package leader_test

import (
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-playground/assert/v2"
	"github.com/redis/go-redis/v9"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func TestNewRedisLeader(t *testing.T) {
	t.Parallel()
	ctx := leader.Ctx

	type TestCase struct {
		description string
		key         string
		wantLeader  *leader.RedisLeader
		wantIsLdr   error
		isLdr       bool
		ldrUUID     string
	}

	testCases := []TestCase{
		{
			description: "Create Leader Failure",
			key:         "",
			wantLeader:  nil,
			wantIsLdr:   errors.New("expected error"),
		},
		{
			description: "Create Leader Success",
			key:         "leader:uuid",
			wantLeader:  nil,
			wantIsLdr:   errors.New("expected error"),
		},
	}

	// Run miniredis server
	miniServer := miniredis.RunT(t)
	defer miniServer.Close()

	// Set up miniredis client no leader
	miniClient := redis.NewClient(&redis.Options{
		Addr:       miniServer.Addr(),
		ClientName: "testclient1",
	})

	_ = miniClient.Set(ctx, "leader:uuid", "key", leader.LeaderTTL)

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			ldr, _ := leader.NewRedisLeader(miniClient, testCase.key)
			assert.Equal(t, ldr, ldr.Key)
		})
	}
}
