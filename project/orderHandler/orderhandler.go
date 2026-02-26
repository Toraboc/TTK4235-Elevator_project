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

	orderHandler.worldView = NewWorldView()

	return &orderHandler
}

func (orderHandler *OrderHandler) GetWorldView() WorldView {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	return orderHandler.worldView.Clone()
}

func (orderHandler *OrderHandler) MergeWorldView(sourceNodeId NodeId, sourceNodeState ElevatorState, sourceOrders Orders) {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	orderHandler.worldView.Merge(sourceNodeId, sourceNodeState, sourceOrders)
}

func (orderHandler *OrderHandler) UpdateConnectedNodes(connectedNodes NodeIdSet) {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	orderHandler.worldView.ConnectedNodes = connectedNodes
}

func (orderHandler *OrderHandler) GetConfirmedOrders() ConfirmedOrders {
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	return orderHandler.worldView.GetConfirmedOrders()
}

func (orderHandler *OrderHandler) GetNextTargetFloor() (int, error){
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	return orderHandler.worldView.GetNextTargetFloor()
}

func (orderHandler *OrderHandler) ChangeElevatorState(state ElevatorState){
	orderHandler.mu.Lock()
	defer orderHandler.mu.Unlock()

	orderHandler.worldView.ElevatorStates[GetMyId()] = state
	orderHandler.worldView.hallRequestAssigner()
}


