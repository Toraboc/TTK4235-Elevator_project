package network

import (
	"fmt"
	"sort"
	"sync"
	"time"

	. "project/orders"
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

type KnownNodes struct {
	mu       sync.Mutex
	LastSeen map[NodeId]time.Time
}

func newNodesAwareOfMe() *NodesAwareOfMe {
	return &NodesAwareOfMe{knowsAboutMe: make(map[NodeId]KnowsAboutMe)}
}

// newKnownNodes creates an initialized KnownNodes.
func newKnownNodes() *KnownNodes {
	return &KnownNodes{LastSeen: make(map[NodeId]time.Time)}
}

// nodeSeen records that the given IP was observed now.
func (nodeSet *KnownNodes) nodeSeen(id NodeId) {
	nodeSet.mu.Lock()
	nodeSet.LastSeen[id] = time.Now()
	nodeSet.mu.Unlock()
}

// updateConnectedNodes prunes stale entries and updates the sorted list of active peer IPs via UpdateConnectedNodes.
func (nodes *KnownNodes) updateConnectedNodes() {
	nodes.mu.Lock()
	defer nodes.mu.Unlock()

	for id, seenAt := range nodes.LastSeen {
		if time.Since(seenAt) > StaleThreshold {
			delete(nodes.LastSeen, id)
		}
	}
	ids := make([]NodeId, 0, len(nodes.LastSeen))
	for id := range nodes.LastSeen {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	//fmt.Println("Active peers:", NodeIdListToStrings(ids)) // DEBUG

	UpdateConnectedNodes(ids)
}

// GetKnowsAboutMe returns a NodeIdSet of nodes that know about me.
func (nodesAwareOfMe *NodesAwareOfMe) KnowsMe() NodeIdSet {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()

	set := make(NodeIdSet)
	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		if entry.Node {
			set.Add(id)
		}
	}
	return set
}

// updateKnowsMe updates the knowsAboutMe based on the received SyncMessage.
func (nodesAwareOfMe *NodesAwareOfMe) update(syncMsg SyncMessage) {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()

	for i := range syncMsg.KnownNodes {
		if syncMsg.KnownNodes[i] == GetMyId() {
			entry := nodesAwareOfMe.knowsAboutMe[syncMsg.Id]
			entry.Node = true
			entry.LastReceived = time.Now()
			nodesAwareOfMe.knowsAboutMe[syncMsg.Id] = entry
		}
	}
}

// purgeStaleKnowsMe marks nodes as not knowing about me if they haven't sent a SyncMessage in a while.
func (nodesAwareOfMe *NodesAwareOfMe) purgeStale() {
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()
	for id, entry := range nodesAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
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
		fmt.Printf("%s: %t, ", NodeIdtoString(id), entry.Node)
	}
	fmt.Println()
}
