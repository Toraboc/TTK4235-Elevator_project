package orderHandler

import (
	. "project/shared"
)

type OrderHandlerInterface struct {
	ConnectedNodesUpdateCh <-chan NodeIdSet
	WorldViewMergeCh       <-chan SyncView
	ElevatorStateCh        <-chan ElevatorState
	OrderCompletedCh       <-chan OrderCompletedEvent
	NewOrderCh             <-chan NewOrderEvent
	WorldViewReqCh         chan chan WorldView
	ConfirmedOrdersCh      chan<- ConfirmedOrders
	TargetFloorCh          chan<- int
}

func RequestWorldView(requestCh chan chan WorldView) WorldView {
	responseCh := make(chan WorldView)
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
