package orderHandler

import (
	"sync"
	. "project/shared"
)

type OrderHandler struct {
	worldView WorldView
	mu        sync.Mutex
}

func NewOrderHandler() *OrderHandler {
	var orderHandler OrderHandler

	orderHandler.worldView = newWorldView()

	return &orderHandler
}

func (orderHandler *OrderHandler) GetWorldView() WorldView {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	return orderHandler.worldView.clone()
}

func (orderHandler *OrderHandler) MergeWorldView(sourceNodeId NodeId, sourceNodeState ElevatorState, sourceOrders Orders) {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	orderHandler.worldView.merge(sourceNodeId, sourceNodeState, sourceOrders)
}

func (orderHandler *OrderHandler) UpdateConnectedNodes(connectedNodes NodeIdSet) {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	orderHandler.worldView.ConnectedNodes = connectedNodes
}

func (orderHandler *OrderHandler) GetConfirmedOrders() ConfirmedOrders {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	return orderHandler.worldView.getConfirmedOrders()
}

func (orderHandler *OrderHandler) GetNextTargetFloor() (int, error){
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	return orderHandler.worldView.getNextTargetFloor()
}

func (orderHandler *OrderHandler) ChangeElevatorState(state ElevatorState){
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	orderHandler.worldView.ElevatorStates[GetMyId()] = state
	orderHandler.worldView.hallRequestAssigner()
}

func (orderHandler *OrderHandler) UpdateNewOrder(floor int, orderType OrderType){
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	myId := GetMyId()
	myOrders := orderHandler.worldView.Orders[myId]

	switch orderType {
	case HALLUP:
		if myOrders.HallUpOrders[floor] == NO_ORDER{
			myOrders.HallUpOrders[floor] = UNCONFIRMED
		}
	case HALLDOWN:
		if myOrders.HallDownOrders[floor] == NO_ORDER{
			myOrders.HallDownOrders[floor] = UNCONFIRMED
		}
	case CAB:
		myCabOrders := myOrders.CabOrders[myId]
		if myCabOrders[floor] ==NO_ORDER{
			myCabOrders[floor] = UNCONFIRMED
		}
		myOrders.CabOrders[myId] = myCabOrders
	default:
		return
	}

	orderHandler.worldView.Orders[myId] = myOrders
}


func (orderHandler *OrderHandler) UpdateFinishedOrder(floor int, orderType OrderType){
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	myId := GetMyId()
	myOrders := orderHandler.worldView.Orders[myId]

	switch orderType {
	case HALLUP:
		myOrders.HallUpOrders[floor] = FINISHED
	case HALLDOWN:
		myOrders.HallDownOrders[floor] = FINISHED
	case CAB:
		myCabOrders := myOrders.CabOrders[myId]
		myCabOrders[floor] = FINISHED
		myOrders.CabOrders[myId] = myCabOrders
	default:
		return
	}
	orderHandler.worldView.Orders[myId] = myOrders
}
