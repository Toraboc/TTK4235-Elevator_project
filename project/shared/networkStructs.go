package shared

import (
	"time"
)

type NetworkNode struct {
	Id       NodeId
	LastSync time.Time
	KnowsMe  bool
}

type SyncMessage struct {
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
	SendTime   time.Time
}
