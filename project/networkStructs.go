package main

import (
	"time"
)

type NetworkNode struct {
    id       NodeId
    lastSync time.Time
    knowsMe  bool
}

type SyncMessage struct {
    orders     Orders
    myState    ElevatorState
    knownNodes []NodeId
}

type NetworkState struct {
    connectedNodes []NetworkNode
}
