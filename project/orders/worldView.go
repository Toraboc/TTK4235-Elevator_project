package orders

import (
	. "project/shared"
)

//TODO: Lage no orderhandler og greier med mutex

func CreateWorldView() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	myId := GetMyId()
	worldView.ConnectedNodes.Add(myId)

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)
	worldView.Orders = make(map[NodeId]Orders)
	worldView.Orders[myId] = CreateOrders(myId)

	return worldView
}

func (wv *WorldView) Clone() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	worldView.ConnectedNodes.Concat(wv.ConnectedNodes)

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)
	for nodeId, elevatorState := range wv.ElevatorStates {
		worldView.ElevatorStates[nodeId] = elevatorState
	}

	worldView.Orders = make(map[NodeId]Orders)
	for nodeId, orders := range wv.Orders {
		worldView.Orders[nodeId] = orders.Clone()
	}

	worldView.AssignedHallUpOrders = wv.AssignedHallUpOrders
	worldView.AssignedHallDownOrders = wv.AssignedHallDownOrders
	worldView.AssignedCabOrders = wv.AssignedCabOrders

	return worldView
}

// This will only sync the orders and elevatorStates
func (worldView *WorldView) Merge(sourceNodeId NodeId, sourceNodeState ElevatorState, sourceOrders Orders) {
	//Oppdaterer staten til heisen m. syncmelding
	worldView.ElevatorStates[sourceNodeId] = sourceNodeState

	// Sync merged orders for source node.
	worldView.Orders[sourceNodeId] = sourceOrders.Clone()
	//TODO: update cyclic counter

	// This must also be called if our own elevatorsstate changes
	worldView.hallRequestAssigner()
}

// This function will receive updates from the elevator
func (worldView *WorldView) ChangeElevatorState(state ElevatorState) {

	worldView.ElevatorStates[GetMyId()] = state
	worldView.hallRequestAssigner()
}
