package worker

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/supertylerc/sre-portfolio/orchestrator/go/task"
	"log/slog"
)

type Worker struct {
	Name      string
	Queue     queue.Queue
	Db        map[uuid.UUID]*task.Task
	TaskCount int
}

func (w *Worker) CollectStats() {
	slog.Debug("I will collect stats")
}

func (w *Worker) RunTask() {
	slog.Debug("I will start or stop a task")
}

func (w *Worker) StartTask() {
	slog.Debug("I will start a task.")
}

func (w *Worker) StopTask() {
	slog.Debug("I will stop a task")
}
