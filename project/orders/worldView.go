package orders

import (
	. "project/shared"
)

func CreateWorldView(nodeId NodeId) Worldview {
	var worldView Worldview

	worldView.ConnectedNodes = make([]NodeId, 1)
	worldView.ConnectedNodes[0] = nodeId

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)

	worldView.Orders.HallUpOrders = CreateOrderList()
	worldView.Orders.HallDownOrders = CreateOrderList()
	worldView.Orders.CabOrders = make(map[NodeId][NumberOfFloors]Order)
	worldView.Orders.CabOrders[nodeId] = CreateOrderList()

	return worldView
}

func GetWorldview(nodeId NodeId) Worldview {
	// TODO: This function should not take this parameter
	// Remember to make this threadsafe

	worldView := CreateWorldView(nodeId)

	worldView.Orders.HallDownOrders[1].LastEvent = NEW
	worldView.Orders.HallDownOrders[2].LastEvent = NEW
	worldView.Orders.HallDownOrders[3].LastEvent = NEW

	return worldView
}
