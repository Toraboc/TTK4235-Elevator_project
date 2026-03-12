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

func (knownNodes *KnownNodes) nodeSeen(id NodeId) {
	knownNodes.mu.Lock()
	knownNodes.LastSeen[id] = time.Now()
	knownNodes.mu.Unlock()
}

func (knownNodes *KnownNodes) pruneStale() bool {
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()
	changed := false

	for id, seenAt := range knownNodes.LastSeen {
		if time.Since(seenAt) > StaleThreshold {
			delete(knownNodes.LastSeen, id)
			changed = true
		}
	}
	return changed
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
