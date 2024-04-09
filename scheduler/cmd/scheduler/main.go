package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	internal_leader "github.com/supertylerc/scheduler/scheduler/internal/leader"
	"github.com/supertylerc/scheduler/scheduler/pkg/leader"
)

func main() {
	internal_leader.ViperConfig()
	internal_leader.LogConfig()

	ctx := context.Background()

	client, err := internal_leader.CreateRedisClient(ctx, 0)
	if err != nil {
		slog.Error("Unable to create Redis Client", slog.String("err", err.Error()))
		os.Exit(1)
	}

	redisLeaderKey := viper.Get("REDIS_LEADER_KEY").(string)

	ldr, err := leader.NewRedisLeader(client, redisLeaderKey)
	if err != nil {
		slog.Error("Unable to create leader", slog.String("err", err.Error()))
		os.Exit(1)
	}

	os.Exit(run(ldr))
}

func run(ldr leader.Leader) int {
	go internal_leader.Metrics()

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
