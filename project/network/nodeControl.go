package network

import (
	"fmt"
	"sync"
	"time"

	. "project/shared"
)

type NetworkNode struct {
	knowsMe  bool
	lastSeen time.Time
}

type NodeControl struct {
	nodes                  map[NodeId]*NetworkNode
	mu                     sync.Mutex
}

func newNodeControl(connectedNodesUpdateCh chan<- NodeIdSet) *NodeControl {
	var nodeControl NodeControl
	nodeControl.nodes = make(map[NodeId]*NetworkNode)

	go nodeControl.updateConnectedNodes(connectedNodesUpdateCh)

	return &nodeControl
}

// Nodes that we have receive a sync message from within the time limit, and that has received one of ours sync messages.
func (nodeControl *NodeControl) getConnectedNodes() NodeIdSet {
	connectedNodes := make(NodeIdSet)
	connectedNodes.Add(GetMyId())

	nodeControl.mu.Lock()
	defer nodeControl.mu.Unlock()

	for nodeId, networkNode := range nodeControl.nodes {
		if networkNode.knowsMe && time.Since(networkNode.lastSeen) < StaleThreshold {
			connectedNodes.Add(nodeId)
		}
	}

	return connectedNodes
}

func (nodeControl *NodeControl) updateConnectedNodes(connectedNodesUpdateCh chan<- NodeIdSet) {
	lastConnectedNodes := nodeControl.getConnectedNodes()
	connectedNodesUpdateCh <- lastConnectedNodes

	ticker := time.NewTicker(time.Second / PruneHz)
	defer ticker.Stop()

	for range ticker.C {
		connectedNodes := nodeControl.getConnectedNodes()

		if !lastConnectedNodes.Equals(connectedNodes) {
			lastConnectedNodes = connectedNodes
			connectedNodesUpdateCh <- connectedNodes.Clone()
			fmt.Printf("Connected nodes: %v\n", connectedNodes)
		}
	}
}

// This function return all the nodes that we have received a sync message from within the time limit.
func (nodeControl *NodeControl) getKnownNodes() NodeIdSet {
	knownNodes := make(NodeIdSet)
	nodeControl.mu.Lock()
	defer nodeControl.mu.Unlock()

	for nodeId, netNetworkNode := range nodeControl.nodes {
		if time.Since(netNetworkNode.lastSeen) < StaleThreshold {
			knownNodes.Add(nodeId)
		}
	}

	return knownNodes
}

func (nodeControl *NodeControl) incommingSync(sourceNodeId NodeId, knownNodes NodeIdSet) {
	nodeControl.mu.Lock()
	defer nodeControl.mu.Unlock()

	networkNode, exists := nodeControl.nodes[sourceNodeId]
	if !exists {
		networkNode = &NetworkNode{}
		nodeControl.nodes[sourceNodeId] = networkNode
	}

	networkNode.lastSeen = time.Now()
	networkNode.knowsMe = knownNodes.Contains(GetMyId())
}

