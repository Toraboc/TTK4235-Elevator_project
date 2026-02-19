package shared

import (
	"fmt"
)

type NodeId [4]byte // Change to int32

type NodeIdSet map[NodeId]struct{}

func (set NodeIdSet) Contains(nodeId NodeId) bool {
	_, exists := set[nodeId]
	return exists
}

func (set NodeIdSet) Add(nodeId NodeId) {
	set[nodeId] = struct{}{}
}

func (set NodeIdSet) Concat(other NodeIdSet) {
	for nodeId := range other {
		set.Add(nodeId)
	}
}

func CreateNodeIdSet(nodeIds []NodeId) NodeIdSet {
	set := make(map[NodeId]struct{})
    for _, node := range nodeIds {
        set[node] = struct{}{}
    }
	return set
}


func (nodeId NodeId) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", nodeId[0], nodeId[1], nodeId[2], nodeId[3])
}