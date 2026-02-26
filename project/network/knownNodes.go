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

// newKnownNodes creates an initialized KnownNodes.
func newKnownNodes() *KnownNodes {
	return &KnownNodes{LastSeen: make(map[NodeId]time.Time)}
}

// nodeSeen records that the given IP was observed now.
func (knownNodes *KnownNodes) nodeSeen(id NodeId) {
	knownNodes.mu.Lock()
	knownNodes.LastSeen[id] = time.Now()
	knownNodes.mu.Unlock()
}

// pruneStale removes nodes that haven't been seen for a while.
func (knownNodes *KnownNodes) pruneStale() {
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()

	for id, seenAt := range knownNodes.LastSeen {
		if time.Since(seenAt) > StaleThreshold {
			delete(knownNodes.LastSeen, id)
		}
	}
}

// Print displays the known nodes and their last seen times.
func (knownNodes *KnownNodes) Print() {
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()

	fmt.Printf("Known nodes: ")
	for id, seenAt := range knownNodes.LastSeen {
		fmt.Printf("%v (last seen: %s), ", id, seenAt.Format(time.RFC3339))
	}
	fmt.Println()
}
