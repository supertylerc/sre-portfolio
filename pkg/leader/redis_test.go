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

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			myClient := miniClient.Set(ctx, "leader:uuid", "key", leader.LeaderTTL)
			ldr, _ := leader.NewRedisLeader(miniClient, testCase.key)
			assert.Equal(t, myClient, ldr.Key)
		})
	}

	// err := miniClient.Set(ctx, "leader:uuid", "key", leader.LeaderTTL).Err()
	// if err != nil {
	//	t.Errorf("Failed to set initial key/value")
	// }

	// miniServer.FastForward(leader.LeaderTTL)

	// Set up leader client on miniClient with RedisLeader Options{}

	// ldr, err := leader.NewRedisLeader(miniClient, "leader:uuid")
	//if err != nil {
	//	t.Errorf("NewRedisLeader() failed %v", err)
	// }

	// slog.Info("New Leader", "contains", ldr.UUID)

	// err = ldr.WriteLeader()
	// if err != nil {
	//	t.Errorf("WriteLeader() failed %v", err)
	// }

	// isLdr, err := ldr.IsCurrentLeader()
	// if err != nil || isLdr == false {
	//	t.Errorf("Call to IsCurrentLeader() failed %v", err)
	// }

	// slog.Info("Current Leader", "status", isLdr)
	// slog.Info("Current Leader", "UUID", ldr.UUID)
}
