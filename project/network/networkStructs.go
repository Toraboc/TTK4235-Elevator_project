package network

import (
	"time"
	. "project/orderHandler"
	. "project/shared"
)

type SyncMessage struct {
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
	SendTime   time.Time
}
