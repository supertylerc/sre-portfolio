package main

import (
"context"
	"log/slog"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"

	"github.com/supertylerc/sre-portfolio/orchestrator/go/internal"
	"github.com/supertylerc/sre-portfolio/orchestrator/go/manager"
	"github.com/supertylerc/sre-portfolio/orchestrator/go/node"
	"github.com/supertylerc/sre-portfolio/orchestrator/go/task"
	"github.com/supertylerc/sre-portfolio/orchestrator/go/worker"
"github.com/sethvargo/go-envconfig"
)

func main() {
ctx := context.Background()
	var c internal.Config
if err := envconfig.Process(ctx, &c); err != nil {
    slog.Debug("Failed to process envconfig", "err", err)
  }

	internal.LogConfig(&c)
	t := task.Task{
		ID:     uuid.New(),
		Name:   "Task-1",
		State:  task.Pending,
		Image:  "Image-1",
		Memory: 1024,
		Disk:   1,
	}
	te := task.TaskEvent{
		ID:        uuid.New(),
		State:     task.Pending,
		Timestamp: time.Now(),
		Task:      t,
	}
	slog.Debug("Task created", "Task", t)
	slog.Debug("TaskEvent created", "TaskEvent", te)

	w := worker.Worker{
		Name:  "Worker-1",
		Queue: *queue.New(),
		Db:    make(map[uuid.UUID]*task.Task),
	}
	slog.Debug("Worker created", "Worker", w)
	w.CollectStats()
	w.RunTask()
	w.StartTask()
	w.StopTask()

	m := manager.Manager{
		Pending: *queue.New(),
		TaskDb:  make(map[string][]*task.Task),
		EventDb: make(map[string][]*task.TaskEvent),
		Workers: []string{w.Name},
	}
	slog.Debug("Manager created", "Manager", m)
	m.SelectWorker()
	m.UpdateTasks()
	m.SendWork()

	n := node.Node{
		Name:   "Node-1",
		Ip:     "192.168.1.1",
		Cores:  4,
		Memory: 1024,
		Disk:   25,
		Role:   "worker",
	}
	slog.Debug("Node created", "Node", n)
}
