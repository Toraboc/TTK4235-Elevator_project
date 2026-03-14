package network

import (
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
		updateConnectedNodes(knownNodes, nodesAwareOfMe, connectedNodesUpdateCh)
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
		updateConnectedNodes(knownNodes, nodesAwareOfMe, connectedNodesUpdateCh)
	}
}
