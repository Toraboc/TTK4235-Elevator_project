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

func (nodesAwareOfMe *NodesAwareOfMe) update(syncMsg SyncMessage, knownNodes *KnownNodes, connectedNodesUpdateCh chan<- NodeIdSet) {
	nodesAwareOfMe.mu.Lock()
	myId := GetMyId()

	changed := false
	for _, nodeId := range syncMsg.KnownNodes {
		if nodeId == myId {
			if entry, _ := nodesAwareOfMe.knowsAboutMe[syncMsg.Id]; !entry.Node {
				changed = true
			}
			entry := nodesAwareOfMe.knowsAboutMe[syncMsg.Id]
			entry.Node = true
			entry.LastReceived = time.Now()
			nodesAwareOfMe.knowsAboutMe[syncMsg.Id] = entry
		}
	}
	nodesAwareOfMe.mu.Unlock()

	if changed {
		nodeUpdate(knownNodes, nodesAwareOfMe, connectedNodesUpdateCh)
	}
}

func (nodesAwareOfMe *NodesAwareOfMe) pruneStale(knownNodes *KnownNodes, connectedNodesUpdateCh chan<- NodeIdSet) {
	nodesAwareOfMe.mu.Lock()

	changed := false
	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
			if entry.Node {
				changed = true
			}
			entry.Node = false
			nodesAwareOfMe.knowsAboutMe[id] = entry
		}
	}
	nodesAwareOfMe.mu.Unlock()

	if changed {
		nodeUpdate(knownNodes, nodesAwareOfMe, connectedNodesUpdateCh)
	}
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
