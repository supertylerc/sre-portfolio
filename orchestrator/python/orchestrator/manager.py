import uuid
import queue

from orchestrator import task


class Manager:
    name: str
    pending: queue.Queue
    task_db: dict[str, list[task.Task]]
    event_db: dict[str, list[task.TaskEvent]]
    workers: list[str]
    worker_task_map: dict[str, list[uuid.UUID]]
    task_worker_map: dict[uuid.UUID, str]

    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            self.__setattr__(k, v)

    def select_worker(self):
        print("I will select a worker")

    def update_tasks(self):
        print("I will update tasks")

    def send_work(self):
        print("I will send work to a worker")
