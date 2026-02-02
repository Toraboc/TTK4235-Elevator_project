package main

import (
	"sync"
	"time"
)

type NetworkNode struct {
	id       NodeId
	lastSync time.Time
	knowsMe  bool
}

type SyncMessage struct {
	id         NodeId
	orders     Orders
	myState    ElevatorState
	knownNodes []NodeId
}

type NetworkState struct {
	connectedNodes []NetworkNode
}

type KnownNodeSet struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
}
