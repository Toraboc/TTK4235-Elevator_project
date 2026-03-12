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

func (nodesAwareOfMe *NodesAwareOfMe) update(syncMsg SyncMessage, nodeUpdateCh chan<- int) {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()
	myId := GetMyId()

	for _, nodeId := range syncMsg.KnownNodes {
		if nodeId == myId {
			if entry, _ := nodesAwareOfMe.knowsAboutMe[syncMsg.Id]; entry.Node == false {
				nodeUpdateCh <- 1
			}
			entry := nodesAwareOfMe.knowsAboutMe[syncMsg.Id]
			entry.Node = true
			entry.LastReceived = time.Now()
			nodesAwareOfMe.knowsAboutMe[syncMsg.Id] = entry
		}
	}
}

func (nodesAwareOfMe *NodesAwareOfMe) pruneStale(nodeUpdateCh chan<- int) {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()

	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
			if entry.Node == true {
				nodeUpdateCh <- 1
			}
			entry.Node = false
			nodesAwareOfMe.knowsAboutMe[id] = entry
		}
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
