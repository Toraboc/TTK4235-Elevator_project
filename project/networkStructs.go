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
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
}

type NetworkState struct {
	connectedNodes []NetworkNode
}

type KnowsMe struct {
	node map[NodeId]bool
	mu   sync.Mutex
}

type KnownNodeSet struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
}
