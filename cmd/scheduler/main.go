package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supertylerc/scheduler/internal"
	"github.com/supertylerc/scheduler/pkg/leader"
)

func main() {
	internal.ViperConfig()
	internal.LogConfig()

	ldr, err := internal.CreateRedisLeader()
	if err != nil {
		slog.Error("Unable to create leader", slog.String("err", err.Error()))
		os.Exit(1)
	}

	os.Exit(run(ldr))
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
