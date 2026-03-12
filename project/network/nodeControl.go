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

func pruneNodes(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe, nodeUpdateCh chan<- int) {
	ticker := time.NewTicker(time.Second / PruneHz)
	defer ticker.Stop()
	for range ticker.C {
		knownNodes.pruneStale(nodeUpdateCh)
		nodesAwareOfMe.pruneStale(nodeUpdateCh)
	}
}

func nodeUpdate(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe, connectedNodesUpdateCh chan<- NodeIdSet, nodeUpdateCh <-chan int) {
	for range nodeUpdateCh {
		connectedNodes := getConnectedNodes(knownNodes, nodesAwareOfMe)
		connectedNodesUpdateCh <- connectedNodes
		fmt.Printf("Connected nodes: %v\n", connectedNodes)
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
