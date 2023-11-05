package leader

import (
	"fmt"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestNewRedisLeader(t *testing.T) {

	// start miniredis instance
	mr := miniredis.RunT(t)

	// Create redis client
	mrClient := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	key := "bGVhZGVyOnV1aWQ="
	ldr, err := NewRedisLeader(mrClient, key)
	if err != nil {
		fmt.Println("Failed to create leader")
	}

	// Print ldr object for now
	fmt.Println(ldr)
}
