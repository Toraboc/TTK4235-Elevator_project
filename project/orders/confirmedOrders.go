package orders

import (
	. "project/shared"
)

type ConfirmedOrders struct {
	HallUp   [NumberOfFloors]bool
	HallDown [NumberOfFloors]bool
	Cab      [NumberOfFloors]bool
}

func findConfirmedOrdersInArray(orders [NumberOfFloors]Order, nodeId NodeId) [NumberOfFloors]bool {
	var confirmed [NumberOfFloors]bool

	for floor := 0; floor < NumberOfFloors; floor++ {
		isConfirmed := true
		status, exists := orders[floor][nodeId]
		if !exists || status != CONFIRMED {
			isConfirmed = false
		}
		confirmed[floor] = isConfirmed
	}
	return confirmed
}

func (worldView *WorldView) GetConfirmedOrders() ConfirmedOrders {
	var confirmedOrders ConfirmedOrders

	confirmedOrders.HallUp = findConfirmedOrdersInArray(worldView.Orders.HallUpOrders, getOwnId())
	confirmedOrders.HallDown = findConfirmedOrdersInArray(worldView.Orders.HallDownOrders, getOwnId())
	confirmedOrders.Cab = make(map[NodeId][NumberOfFloors]bool)
	confirmedOrders.Cab = findConfirmedOrdersInArray(worldView.Orders.cabOrders[getOwnId()], getOwnId())

	return confirmedOrders
}
