package orderHandler

import . "project/shared"

type NewOrderEvent struct {
	Floor     int
	OrderType OrderType
}

type WorldViewRequestCh chan chan WorldView

type OrderChannels struct {
	ConnectedNodesUpdateCh chan NodeIdSet
	WorldViewMergeCh       chan SyncView
	ElevatorStateCh        chan ElevatorState
	OrderCompletedCh       chan OrderCompleted
	NewOrderCh             chan NewOrderEvent
	WorldViewReqCh         WorldViewRequestCh
	ConfirmedOrdersCh      chan ConfirmedOrders
	TargetFloorCh          chan int
}

func NewOrderChannels() OrderChannels {
	return OrderChannels{
		ConnectedNodesUpdateCh: make(chan NodeIdSet, 1),
		WorldViewMergeCh:       make(chan SyncView, 1),
		ElevatorStateCh:        make(chan ElevatorState, 1),
		OrderCompletedCh:       make(chan OrderCompleted, 10),
		NewOrderCh:             make(chan NewOrderEvent, 10),
		WorldViewReqCh:         make(WorldViewRequestCh, 1),
		ConfirmedOrdersCh:      make(chan ConfirmedOrders, 1),
		TargetFloorCh:          make(chan int, 1),
	}
}

func RequestWorldView(requestCh WorldViewRequestCh) WorldView {
	responseCh := make(chan WorldView)
	requestCh <- responseCh
	return <-responseCh
}

