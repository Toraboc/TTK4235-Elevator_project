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

// newNodesAwareOfMe creates an initialized NodesAwareOfMe.
func newNodesAwareOfMe() *NodesAwareOfMe {
	return &NodesAwareOfMe{knowsAboutMe: make(map[NodeId]KnowsAboutMe)}
}

// updateKnowsMe updates the knowsAboutMe based on the received SyncMessage.
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

// pruneStale marks nodes as not knowing about me if they haven't sent a SyncMessage in a while.
func (nodesAwareOfMe *NodesAwareOfMe) pruneStale() {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()
	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
			entry.Node = false
			nodesAwareOfMe.knowsAboutMe[id] = entry
		}
	}
}

// Print displays the nodes that know about me.
func (nodesAwareOfMe *NodesAwareOfMe) Print() {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()

	fmt.Printf("Knows about me: ")
	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		fmt.Printf("%v: %t, ", id, entry.Node)
	}
	fmt.Println()
}
