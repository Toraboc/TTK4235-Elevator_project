package orderHandler

import (
	. "project/shared"
)

type OrderHandlerInterface struct {
	ConnectedNodesUpdateCh <-chan NodeIdSet
	SyncMergeCh            <-chan SyncData
	ElevatorStateCh        <-chan ElevatorState
	OrderCompletedCh       <-chan OrderCompletedEvent
	NewOrderCh             <-chan NewOrderEvent
	RequestSyncCh          chan chan SyncData
	ConfirmedOrdersCh      chan<- ConfirmedOrders
	TargetFloorCh          chan<- int
}

type SyncData struct {
	NodeId        NodeId
	ElevatorState ElevatorState
	Orders        Orders
}

func RequestSyncView(requestCh chan chan SyncData) SyncData {
	responseCh := make(chan SyncData)
	requestCh <- responseCh
	return <-responseCh
}

type OrderCompletedEvent struct {
	Floor     int
	Direction Direction
}

type NewOrderEvent struct {
	Floor     int
	OrderType OrderType
}
