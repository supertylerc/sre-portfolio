from enum import Enum
import uuid


class State(Enum):
    Pending = 0
    Scheduled = 1
    Running = 2
    Completed = 3
    Failed = 4


class Task:
    id: uuid.UUID
    container_id: str
    name: str
    state: State
    image: str
    cpu: float
    memory: int
    disk: int
    # TC: TODO: Figure out how to do this in Python Docker lib
    # exposed_ports: nat.PortSet
    port_bindings: dict[str, str]
    restart_policy: str
    start_time: float
    finish_time: float

    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            self.__setattr__(k, v)


class TaskEvent:
    id: uuid.UUID
    state: State
    timestamp: float
    task: Task

    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            self.__setattr__(k, v)
