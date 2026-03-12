package orderHandler

import . "project/shared"

type NewOrderEvent struct {
	Floor     int
	OrderType OrderType
}

type WorldViewRequestCh chan chan WorldView
type ConfirmedOrdersRequestCh chan chan ConfirmedOrders

type OrderChannels struct {
	ConnectedNodesUpdateCh chan NodeIdSet
	WorldViewMergeCh       chan SyncView
	ElevatorStateCh        chan ElevatorState
	OrderCompletedCh       chan OrderCompleted
	NewOrderCh             chan NewOrderEvent
	WorldViewReqCh         WorldViewRequestCh
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
		WorldViewReqCh:         make(WorldViewRequestCh),
		ConfirmedOrdersReqCh:   make(ConfirmedOrdersRequestCh),
		TargetFloorCh:          make(chan int),
	}
}

func RequestWorldView(requestCh WorldViewRequestCh) WorldView {
	responseCh := make(chan WorldView)
	requestCh <- responseCh
	return <-responseCh
}

func RequestConfirmedOrders(requestCh ConfirmedOrdersRequestCh) ConfirmedOrders {
	responseCh := make(chan ConfirmedOrders)
	requestCh <- responseCh
	return <-responseCh
}
