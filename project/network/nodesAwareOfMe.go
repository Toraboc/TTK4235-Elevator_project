package network

import (
	"fmt"
	"sync"
	"time"

	. "project/shared"
)

type KnowsAboutMe struct {
	Node         bool
	LastReceived time.Time
}

type NodesAwareOfMe struct {
	mu           sync.Mutex
	knowsAboutMe map[NodeId]KnowsAboutMe
}

func newNodesAwareOfMe() *NodesAwareOfMe {
	return &NodesAwareOfMe{knowsAboutMe: make(map[NodeId]KnowsAboutMe)}
}

func (nodesAwareOfMe *NodesAwareOfMe) update(syncMsg SyncMessage) {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()
	myId := GetMyId()

	for _, nodeId := range syncMsg.KnownNodes {
		if nodeId == myId {
			entry := nodesAwareOfMe.knowsAboutMe[syncMsg.Id]
			entry.Node = true
			entry.LastReceived = time.Now()
			nodesAwareOfMe.knowsAboutMe[syncMsg.Id] = entry
		}
	}
}

func (nodesAwareOfMe *NodesAwareOfMe) pruneStale() bool {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()
	changed := false

	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
			entry.Node = false
			nodesAwareOfMe.knowsAboutMe[id] = entry
			changed = true
		}
	}
	return changed
}

func (nodesAwareOfMe *NodesAwareOfMe) Print() {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()

	fmt.Printf("Knows about me: ")
	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		fmt.Printf("%v: %t, ", id, entry.Node)
	}
	fmt.Println()
}
