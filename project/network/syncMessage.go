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

func createOutgoingSync(orderHandler *OrderHandler, knownNodes *KnownNodes) SyncMessage {
	worldview := orderHandler.GetWorldView()

	syncMsg := SyncMessage{}
	syncMsg.Id = GetMyId()
	syncMsg.Orders = *worldview.Orders[syncMsg.Id].Clone()
	syncMsg.MyState = worldview.ElevatorStates[syncMsg.Id]
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()
	syncMsg.KnownNodes = make([]NodeId, 0, len(knownNodes.LastSeen))
	for id := range knownNodes.LastSeen {
		syncMsg.KnownNodes = append(syncMsg.KnownNodes, id)
	}
	syncMsg.SendTime = time.Now()
	return syncMsg
}
