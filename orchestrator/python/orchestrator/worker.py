import uuid
import queue

from orchestrator import task


class Worker:
    name: str
    queue: queue.Queue
    db: dict[uuid.UUID, task.Task]
    task_count: int

    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            self.__setattr__(k, v)

    def collect_stats(self):
        print("I will collect stats")

    def run_task(self):
        print("I will start or stop a task")

    def start_task(self):
        print("I will start a task.")

    def stop_task(self):
        print("I will stop a task")
