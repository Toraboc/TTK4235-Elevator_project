package network

import (
	"fmt"
	"sync"
	"time"

	. "project/shared"
)

type KnownNodes struct {
	mu       sync.Mutex
	LastSeen map[NodeId]time.Time
}

func newKnownNodes() *KnownNodes {
	return &KnownNodes{LastSeen: make(map[NodeId]time.Time)}
}

func (knownNodes *KnownNodes) nodeSeen(id NodeId, nodesAwareOfMe *NodesAwareOfMe, connectedNodesUpdateCh chan<- NodeIdSet) {
	knownNodes.mu.Lock()
	_, exists := knownNodes.LastSeen[id]
	knownNodes.LastSeen[id] = time.Now()
	knownNodes.mu.Unlock()

	if !exists {
		nodeUpdate(knownNodes, nodesAwareOfMe, connectedNodesUpdateCh)
	}
}

func (knownNodes *KnownNodes) pruneStale(nodesAwareOfMe *NodesAwareOfMe, connectedNodesUpdateCh chan<- NodeIdSet) {
	knownNodes.mu.Lock()

	changed := false
	for id, seenAt := range knownNodes.LastSeen {
		if time.Since(seenAt) > StaleThreshold {
			delete(knownNodes.LastSeen, id)
			changed = true
		}
	}
	knownNodes.mu.Unlock()

	if changed {
		nodeUpdate(knownNodes, nodesAwareOfMe, connectedNodesUpdateCh)
	}
}

func (knownNodes *KnownNodes) print() {
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()

	fmt.Printf("Known nodes: ")
	for id, seenAt := range knownNodes.LastSeen {
		fmt.Printf("%v (last seen: %s), ", id, seenAt.Format(time.RFC3339))
	}
	fmt.Println()
}
