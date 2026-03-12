package orderHandler

import . "project/shared"

type NewOrderEvent struct {
	Floor     int
	OrderType OrderType
}

type ConfirmedOrdersRequestCh chan chan ConfirmedOrders

type OrderChannels struct {
	ConnectedNodesUpdateCh chan NodeIdSet
	WorldViewMergeCh       chan SyncView
	ElevatorStateCh        chan ElevatorState
	OrderCompletedCh       chan OrderCompleted
	NewOrderCh             chan NewOrderEvent
	WorldViewReqCh         chan chan WorldView
	ConfirmedOrdersReqCh   ConfirmedOrdersRequestCh
	TargetFloorCh          chan int
}

func NewOrderChannels() OrderChannels {
	return OrderChannels{
		ConnectedNodesUpdateCh: make(chan NodeIdSet),
		WorldViewMergeCh:       make(chan SyncView),
		ElevatorStateCh:        make(chan ElevatorState),
		OrderCompletedCh:       make(chan OrderCompleted),
		NewOrderCh:             make(chan NewOrderEvent),
		WorldViewReqCh:         make(chan chan WorldView),
		ConfirmedOrdersReqCh:   make(ConfirmedOrdersRequestCh),
		TargetFloorCh:          make(chan int),
	}
}

func (channels OrderChannels) RequestWorldView() WorldView {
	responseCh := make(chan WorldView)
	channels.WorldViewReqCh <- responseCh
	return <-responseCh
}
