package orders

import (
	. "project/shared"
)

type ConfirmedOrders struct {
	HallUp   [NumberOfFloors]OrderStatus
	HallDown [NumberOfFloors]OrderStatus
	Cab      map[NodeId][NumberOfFloors]OrderStatus
}

func isConfirmedByEveryone(nodes []NodeId, connectedNodes []NodeId) bool {
	set := CreateNodeIdSet(nodes) // TODO: I think we need to use this set in many of our datatypes

	for _, connectedNode := range connectedNodes {
		if !set.Contains(connectedNode) {
			return false
		}
	}

	return true
}

func findConfirmedOrdersInArray(orders [NumberOfFloors]Order, connectedNodes []NodeId) [NumberOfFloors]OrderStatus {
	var orderStatus [NumberOfFloors]OrderStatus

	for i := range NumberOfFloors {
		if orders[i].LastEvent == NEW && isConfirmedByEveryone(orders[i].ConfirmedBy, connectedNodes) {
			orderStatus[i] = NEW
		} else {
			orderStatus[i] = COMPLETED
		}
	}

	return orderStatus
}

func GetConfirmedOrders() ConfirmedOrders {
	// TODO: This function should not take nodeId as an argument
	worldview := GetWorldView()

	connectedNodes := worldview.ConnectedNodes
	orders := worldview.Orders
	var confirmedOrders ConfirmedOrders

	confirmedOrders.HallUp = findConfirmedOrdersInArray(orders.HallUpOrders, connectedNodes)
	confirmedOrders.HallDown = findConfirmedOrdersInArray(orders.HallDownOrders, connectedNodes)
	confirmedOrders.Cab = make(map[NodeId][NumberOfFloors]OrderStatus)

	for nodeId, orders := range orders.CabOrders {
		confirmedOrders.Cab[nodeId] = findConfirmedOrdersInArray(orders, connectedNodes)
	}

	return confirmedOrders
}
