package network

import (
	"fmt"
	"time"

	. "project/orderHandler"
	. "project/shared"
)

func getConnectedNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) NodeIdSet {
	set := make(NodeIdSet)
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()

	for id := range knownNodes.LastSeen {
		nodesAwareOfMe.mu.Lock()
		if entry, exists := nodesAwareOfMe.knowsAboutMe[id]; exists && entry.Node {
			set.Add(id)
		}
		nodesAwareOfMe.mu.Unlock()
	}
	return set
}

func pruneNodes(orderHandler *OrderHandler, knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
	ticker := time.NewTicker(time.Second / PruneHz)
	defer ticker.Stop()
	for range ticker.C {
		knownNodes.pruneStale()
		nodesAwareOfMe.pruneStale()
		orderHandler.UpdateConnectedNodes(getConnectedNodes(knownNodes, nodesAwareOfMe))
	}
}

func printConnectedNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
	connectedNodes := getConnectedNodes(knownNodes, nodesAwareOfMe)
	for {
		time.Sleep(1 * time.Second / PrintHz)
		connectedNodes = getConnectedNodes(knownNodes, nodesAwareOfMe)
		fmt.Printf("Connected nodes: ")
		for id := range connectedNodes {
			fmt.Printf("%v, ", id)
		}
		fmt.Println()
	}
}
