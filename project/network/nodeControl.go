package network

import (
	"fmt"
	"time"

	. "project/orderHandler"
	. "project/shared"
)

// GetConnectedNodes returns a NodeIdSet of the nodes that have 2-way communication
func GetConnectedNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) NodeIdSet {
	set := make(NodeIdSet)
	knownNodes.mu.Lock()
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()
	defer knownNodes.mu.Unlock()

	for id := range knownNodes.LastSeen {
		if entry, exists := nodesAwareOfMe.knowsAboutMe[id]; exists && entry.Node {
			set.Add(id)
		}
	}
	return set
}

// pruneNodes periodically prunes stale nodes from knownNodes and nodesAwareOfMe, and updates the connected nodes.
func pruneNodes(orderHandler *OrderHandler, knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
	ticker := time.NewTicker(time.Second / PruneHz) // last number controls how often inactive peers are pruned (Hz)
	defer ticker.Stop()
	for range ticker.C {
		knownNodes.pruneStale()
		nodesAwareOfMe.pruneStale()
		orderHandler.UpdateConnectedNodes(GetConnectedNodes(knownNodes, nodesAwareOfMe))
	}
}

func printConnectedNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
	connectedNodes := GetConnectedNodes(knownNodes, nodesAwareOfMe)
	for {
		time.Sleep(1 * time.Second)
		connectedNodes = GetConnectedNodes(knownNodes, nodesAwareOfMe)
		fmt.Printf("Connected nodes: ")
		for id := range connectedNodes {
			fmt.Printf("%v, ", id)
		}
		fmt.Println()
	}
}
