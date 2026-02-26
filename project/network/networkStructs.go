package network

import (
	"time"
	. "project/orderHandler"
	. "project/shared"
)

type NetworkNode struct {
	Id       NodeId
	LastSync time.Time
	KnowsMe  bool
}

//Dette er kanskje dårlig struktur hilsen Paulius. Unødvendig coupling
type SyncMessage struct {
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
	SendTime   time.Time
}
