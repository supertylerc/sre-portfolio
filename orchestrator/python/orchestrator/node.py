class Node:
    name: str
    ip: str
    cores: int
    memory: int
    memory_allocated: int
    disk: int
    disk_allocated: int
    role: str
    task_count: int

    def __init__(self, **kwargs):
        for k, v in kwargs.items():
            self.__setattr__(k, v)
