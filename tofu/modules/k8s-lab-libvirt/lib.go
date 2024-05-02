package lib

import (
	"fmt"
)

type Node struct {
	Kind   string
	Num    int
	CPU    int `tf:"cpu"`
	Memory int
	Disk   int
}

func ExpandNodes(nodes []Node) map[string]Node {
	nodeMap := make(map[string]Node)
	for _, n := range nodes {
		for i := 0; i < n.Num; i++ {
			k := fmt.Sprintf("%s-%d", n.Kind, i)
			nodeMap[k] = Node{
				Kind:   n.Kind,
				CPU:    n.CPU,
				Memory: n.Memory,
				Disk:   n.Disk,
				Num:    i,
			}
		}
	}
	return nodeMap
}
