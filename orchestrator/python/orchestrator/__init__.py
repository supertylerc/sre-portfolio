import queue
import time
import uuid

from orchestrator import task
from orchestrator import worker
from orchestrator import manager
from orchestrator import node


def main():
    t = task.Task(
        id=uuid.uuid4(),
        name="Task-1",
        state=task.State.Pending,
        image="Image-1",
        memory=1024,
        disk=1,
    )
    te = task.TaskEvent(
        id=uuid.uuid4(),
        state=task.State.Pending,
        timestamp=time.time(),
        task=t,
    )
    print("Task created", "Task", t)
    print("TaskEvent created", "TaskEvent", te)

    w = worker.Worker(
        name="Worker-1",
        queue=queue.Queue(),
        Db={},
    )
    print("Worker created", "Worker", w)
    w.collect_stats()
    w.run_task()
    w.start_task()
    w.stop_task()

    m = manager.Manager(
        pending=queue.Queue(),
        task_db={},
        event_db={},
        workers=[w.name],
    )
    print("Manager created", "Manager", m)
    m.select_worker()
    m.update_tasks()
    m.send_work()

    n = node.Node(
        name="Node-1",
        ip="192.168.1.1",
        cores=4,
        memory=1024,
        disk=25,
        role="worker",
    )
    print("Node created", "Node", n)
