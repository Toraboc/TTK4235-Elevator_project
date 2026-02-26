package orderHandler

import (
	. "project/shared"
)


type WorldView struct {
	Orders                 map[NodeId]Orders
	ConnectedNodes         NodeIdSet
	ElevatorStates         map[NodeId]ElevatorState
	AssignedHallUpOrders   [NumberOfFloors]bool
	AssignedHallDownOrders [NumberOfFloors]bool
	AssignedCabOrders      [NumberOfFloors]bool
}


//TODO: Lage no orderhandler og greier med mutex

func newWorldView() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	myId := GetMyId()
	worldView.ConnectedNodes.Add(myId)

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)
	worldView.Orders = make(map[NodeId]Orders)
	worldView.Orders[myId] = newOrders(myId)

	return worldView
}

func (worldView *WorldView) clone() WorldView {
	var clone WorldView

	clone.ConnectedNodes = make(NodeIdSet)
	clone.ConnectedNodes.Concat(worldView.ConnectedNodes)

	clone.ElevatorStates = make(map[NodeId]ElevatorState)
	for nodeId, elevatorState := range worldView.ElevatorStates {
		clone.ElevatorStates[nodeId] = elevatorState
	}

	clone.Orders = make(map[NodeId]Orders)
	for nodeId, orders := range worldView.Orders {
		clone.Orders[nodeId] = orders.clone()
	}

	clone.AssignedHallUpOrders = worldView.AssignedHallUpOrders
	clone.AssignedHallDownOrders = worldView.AssignedHallDownOrders
	clone.AssignedCabOrders = worldView.AssignedCabOrders

	return clone
}

// This will only sync the orders and elevatorStates
func (worldView *WorldView) merge(sourceNodeId NodeId, sourceNodeState ElevatorState, sourceOrders Orders) {
	//Oppdaterer staten til heisen m. syncmelding
	worldView.ElevatorStates[sourceNodeId] = sourceNodeState

	// Sync merged orders for source node.
	worldView.Orders[sourceNodeId] = sourceOrders.clone()

	worldView.updateCyclicCounter()

	// This must also be called if our own elevatorsstate changes
	worldView.hallRequestAssigner()
}


func (worldView *WorldView) updateCyclicCounter() {
	myId := GetMyId()
	connectedNodes := worldView.ConnectedNodes
	connectedNodes.Remove(myId)

	getHallDown := func (orders Orders) *OrderList {
		return &orders.HallDownOrders
	}
	updateCyclicCounter(worldView.Orders, myId, connectedNodes, getHallDown)

	getHallUp := func (orders Orders) *OrderList {
		return &orders.HallUpOrders
	}
	updateCyclicCounter(worldView.Orders, myId, connectedNodes, getHallUp)

	for nodeId, _ := range connectedNodes {
		getCabOrder := func (orders Orders) *OrderList {
			// TODO: This seems sus, maybe this field selector doesn't need to return a pointer?
			orderList := orders.CabOrders[nodeId]
			return &orderList
		}
		updateCyclicCounter(worldView.Orders, myId, connectedNodes, getCabOrder)
	}

	getMyCab := func (orders Orders) *OrderList {
		orderList := orders.CabOrders[myId]
		return &orderList
	}
	updateCyclicCounter(worldView.Orders, myId, connectedNodes, getMyCab)
}

