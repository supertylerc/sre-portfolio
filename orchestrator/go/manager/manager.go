package manager

import (
	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/supertylerc/sre-portfolio/orchestrator/go/task"
	"log/slog"
)

type Manager struct {
	Name          string
	Pending       queue.Queue
	TaskDb        map[string][]*task.Task
	EventDb       map[string][]*task.TaskEvent
	Workers       []string
	WorkerTaskMap map[string][]uuid.UUID
	TaskWorkerMap map[uuid.UUID]string
}

func (m *Manager) SelectWorker() {
	slog.Debug("I will select a worker")
}

func (m *Manager) UpdateTasks() {
	slog.Debug("I will update tasks")
}

func (m *Manager) SendWork() {
	slog.Debug("I will send work to a worker")
}
