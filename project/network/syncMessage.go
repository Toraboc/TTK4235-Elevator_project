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
	SendTime   time.Time // TODO: Is this needed ?
}

func createOutgoingSync(channels OrderChannels, knownNodes *KnownNodes) SyncMessage {
	worldview := channels.RequestWorldView()

	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()

	syncMsg := SyncMessage{}

	syncMsg.Id = GetMyId()
	syncMsg.Orders = *worldview.Orders[syncMsg.Id].Clone()
	syncMsg.MyState = worldview.ElevatorStates[syncMsg.Id]
	syncMsg.KnownNodes = make([]NodeId, 0, len(knownNodes.LastSeen))
	for id := range knownNodes.LastSeen {
		syncMsg.KnownNodes = append(syncMsg.KnownNodes, id)
	}
	syncMsg.SendTime = time.Now()
	return syncMsg
}
