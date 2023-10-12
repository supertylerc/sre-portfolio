package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"github.com/supertylerc/scheduler/internal"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func main() {
	internal.ViperConfig()
	internal.LogConfig()

	ldr, err := createRedisLeader()
	if err != nil {
		slog.Error("Unable to create leader", slog.String("err", err.Error()))
		os.Exit(1)
	}

	os.Exit(run(ldr))
}

func createRedisLeader() (*leader.RedisLeader, error) {
	redisHost := viper.Get("REDIS_HOST").(string)
	redisPort := viper.Get("REDIS_PORT").(string)
	redisPassword := viper.Get("REDIS_PASSWORD").(string)
	redisLeaderKey := viper.Get("REDIS_LEADER_KEY").(string)
	slog.Debug(
		"Creating leader",
		slog.String("redisHost", redisHost),
		slog.String("redisPort", redisPort),
		slog.String("redisPassword", redisPassword),
		slog.String("redisLeaderKey", redisLeaderKey),
	)

	ldr, err := leader.NewRedisLeader(
		fmt.Sprintf("%s:%s", redisHost, redisPort),
		redisPassword,
		redisLeaderKey,
	)
	if err != nil {
		return &leader.RedisLeader{}, fmt.Errorf("error creating a Redis Leader: %w", err)
	}

	slog.Info(
		"Created new leader",
		slog.String("uuid", ldr.UUID.String()),
		slog.String("key", ldr.Key),
	)

	return ldr, nil
}

func run(ldr leader.Leader) int {
	go internal.Metrics()

	ticker := time.NewTicker(leader.CheckInterval * time.Millisecond)
	done := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				slog.Debug("Writing leader")

				err := ldr.WriteLeader()
				if err != nil {
					slog.Error("Error writing leader", slog.String("err", err.Error()))
				}
			}
		}
	}()

	for {
		<-term
		slog.Debug("Ending Goroutine.")
		done <- struct{}{}

		slog.Debug("Closing gracefully.")

		return 0
	}
}
