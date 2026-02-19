package network

import (
	"fmt"
	"sort"
	"sync"
	"time"

	//. "project/orders"
	. "project/shared"
)

type KnownNodes struct {
	mu       sync.Mutex
	LastSeen map[NodeId]time.Time
}

type KnowsAboutMe struct {
	Node         bool
	LastReceived time.Time
}

type NodesAwareOfMe struct {
	mu           sync.Mutex
	knowsAboutMe map[NodeId]KnowsAboutMe
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
	ids := make([]NodeId, 0, len(knownNodes.LastSeen))
	for id := range knownNodes.LastSeen {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
}

// Print displays the known nodes and their last seen times.
func (knownNodes *KnownNodes) Print() {
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()

	fmt.Printf("Known nodes: ")
	for id, seenAt := range knownNodes.LastSeen {
		fmt.Printf("%s (last seen: %s), ", NodeIdtoString(id), seenAt.Format(time.RFC3339))
	}
	fmt.Println()
}

//________________________________________________________________________________________________________

// newNodesAwareOfMe creates an initialized NodesAwareOfMe.
func newNodesAwareOfMe() *NodesAwareOfMe {
	return &NodesAwareOfMe{knowsAboutMe: make(map[NodeId]KnowsAboutMe)}
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
		fmt.Printf("%s: %t, ", NodeIdtoString(id), entry.Node)
	}
	fmt.Println()
}

//__________________________________________________________________________________________________________

// GetConnectedNodes returns a NodeIdSet of the nodes that have 2-way communication
func GetConnectedNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) NodeIdSet {
	set := make(NodeIdSet)
	knownNodes.mu.Lock()
	nodesAwareOfMe.mu.Lock()
	defer knownNodes.mu.Unlock()
	defer nodesAwareOfMe.mu.Unlock()

	for id := range knownNodes.LastSeen {
		if entry, exists := nodesAwareOfMe.knowsAboutMe[id]; exists && entry.Node {
			set.Add(id)
		}
	}
	return set
}
