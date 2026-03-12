package network

import (
	"fmt"
	"time"

	. "project/shared"
)

func getConnectedNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) NodeIdSet {
	set := make(NodeIdSet)
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()
	nodesAwareOfMe.mu.Lock()
	defer nodesAwareOfMe.mu.Unlock()

	for id := range knownNodes.LastSeen {
		if entry, exists := nodesAwareOfMe.knowsAboutMe[id]; exists && entry.Node {
			set.Add(id)
		}
	}
	return set
}

func pruneNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe, connectedNodesUpdateChannel chan NodeIdSet) {
	ticker := time.NewTicker(time.Second / PruneHz)
	defer ticker.Stop()
	for range ticker.C {
		knownNodes.pruneStale()
		nodesAwareOfMe.pruneStale()
		connectedNodesUpdateChannel <- getConnectedNodes(knownNodes, nodesAwareOfMe)
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
