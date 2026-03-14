package orderHandler

import (
	. "project/shared"
)

type OrderChannels struct {
	ConnectedNodesUpdateCh chan NodeIdSet
	WorldViewMergeCh       chan SyncView
	ElevatorStateCh        chan ElevatorState
	OrderCompletedCh       chan OrderCompletedEvent
	NewOrderCh             chan NewOrderEvent
	WorldViewReqCh         chan chan WorldView
	ConfirmedOrdersCh      chan ConfirmedOrders
	TargetFloorCh          chan int
}

func NewOrderChannels() OrderChannels {
	return OrderChannels{
		ConnectedNodesUpdateCh: make(chan NodeIdSet, 1),
		WorldViewMergeCh:       make(chan SyncView, 1),
		ElevatorStateCh:        make(chan ElevatorState, 1),
		OrderCompletedCh:       make(chan OrderCompletedEvent, 10),
		NewOrderCh:             make(chan NewOrderEvent, 10),
		WorldViewReqCh:         make(chan chan WorldView),
		ConfirmedOrdersCh:      make(chan ConfirmedOrders, 1),
		TargetFloorCh:          make(chan int, 1),
	}
}

func RequestWorldView(requestCh chan chan WorldView) WorldView {
	responseCh := make(chan WorldView)
	requestCh <- responseCh
	return <-responseCh
}

type OrderCompletedEvent struct {
	Floor int
	Direction Direction
}

type NewOrderEvent struct {
	Floor int
	OrderType OrderType
}
