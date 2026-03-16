package network

import (
	. "project/orderHandler"
	. "project/shared"
)

type SyncMessage struct {
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
}

func nodeIdSetToList(nodeIdSet NodeIdSet) []NodeId {
	var nodeIdList []NodeId

	for nodeId := range nodeIdSet {
		nodeIdList = append(nodeIdList, nodeId)
	}

	return nodeIdList
}

func createOutgoingSync(worldViewReqCh chan chan WorldView, nodeControl *NodeControl) SyncMessage {
	worldview := RequestWorldView(worldViewReqCh)

	syncMsg := SyncMessage{}

	syncMsg.Id = GetMyId()
	syncMsg.Orders = *worldview.Orders[syncMsg.Id].Clone()
	syncMsg.MyState = worldview.ElevatorStates[syncMsg.Id]
	syncMsg.KnownNodes = nodeIdSetToList(nodeControl.getKnownNodes())
	return syncMsg
}
