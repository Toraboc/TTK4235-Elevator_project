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

func createOutgoingSync(requestSyncCh chan chan SyncData, nodeControl *NodeControl) SyncMessage {
	syncData := RequestSyncView(requestSyncCh)

	syncMsg := SyncMessage{}

	syncMsg.Id = GetMyId()
	syncMsg.Orders = syncData.Orders
	syncMsg.MyState = syncData.ElevatorState
	syncMsg.KnownNodes = nodeIdSetToList(nodeControl.getKnownNodes())
	return syncMsg
}
