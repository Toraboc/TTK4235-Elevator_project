package shared

import (
	"time"
)

type SyncMessage struct {
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
	SendTime   time.Time
}
