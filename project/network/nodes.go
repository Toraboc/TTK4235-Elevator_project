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

type PeersAwareOfMe struct {
	mu           sync.Mutex
	knowsAboutMe map[NodeId]KnowsAboutMe
}

type KnownNodes struct {
	Mu       sync.Mutex
	LastSeen map[NodeId]time.Time
}

// newKnownNodes creates an initialized KnownNodes.
func newKnownNodes() *KnownNodes {
	return &KnownNodes{LastSeen: make(map[NodeId]time.Time)}
}

// seen records that the given IP was observed now.
func (nodeSet *KnownNodes) nodeSeen(id NodeId) {
	nodeSet.Mu.Lock()
	nodeSet.LastSeen[id] = time.Now()
	nodeSet.Mu.Unlock()
}

// list returns the sorted list of active peer IPs and prunes stale entries.
func (nodes *KnownNodes) listActivePeers() {
	nodes.Mu.Lock()
	defer nodes.Mu.Unlock()

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

func GetKnowsAboutMe(peersAwareOfMe *PeersAwareOfMe) NodeIdSet {
	peersAwareOfMe.mu.Lock()
	defer peersAwareOfMe.mu.Unlock()

	set := make(NodeIdSet)
	for id, entry := range peersAwareOfMe.knowsAboutMe {
		if entry.Node {
			set.Add(id)
		}
	}
	return set
}

// updateKnowsMe updates the knowsAboutMe based on the received SyncMessage.
func updateKnowsMe(syncMsg SyncMessage, otherNodes *PeersAwareOfMe) {
	otherNodes.mu.Lock()
	defer otherNodes.mu.Unlock()

	for i := range syncMsg.KnownNodes {
		if syncMsg.KnownNodes[i] == myId {
			entry := otherNodes.knowsAboutMe[syncMsg.Id]
			entry.Node = true
			entry.LastReceived = time.Now()
			otherNodes.knowsAboutMe[syncMsg.Id] = entry
		}
	}
}

// purgeStaleKnowsMe marks nodes as not knowing about me if they haven't sent a SyncMessage in a while.
func purgeStaleKnowsMe(peersAwareOfMe *PeersAwareOfMe) {
	peersAwareOfMe.mu.Lock()
	defer peersAwareOfMe.mu.Unlock()
	for id, entry := range peersAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
			entry.Node = false
			peersAwareOfMe.knowsAboutMe[id] = entry
		}
	}
}

func printKnowsAboutMe(peersAwareOfMe *PeersAwareOfMe) {
	peersAwareOfMe.mu.Lock()
	defer peersAwareOfMe.mu.Unlock()

	fmt.Printf("Knows about me: ")
	for id, entry := range peersAwareOfMe.knowsAboutMe {
		fmt.Printf("%s: %t, ", NodeIdtoString(id), entry.Node)
	}
	fmt.Println()
}
