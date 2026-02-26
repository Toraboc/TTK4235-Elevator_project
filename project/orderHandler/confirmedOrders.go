package orderHandler

import (
	. "project/shared"
)

type ConfirmedOrders struct {
	HallUp   [NumberOfFloors]bool
	HallDown [NumberOfFloors]bool
	Cab      [NumberOfFloors]bool
}

func findConfirmedOrdersInArray(orders OrderList) [NumberOfFloors]bool {
	var confirmed [NumberOfFloors]bool

	for floor := 0; floor < NumberOfFloors; floor++ {
		confirmed[floor] = orders[floor] == CONFIRMED
	}
	return confirmed
}

func (worldView *WorldView) GetConfirmedOrders() ConfirmedOrders {
	var confirmedOrders ConfirmedOrders
	ownId := GetMyId()

	orders, exists := worldView.Orders[ownId]
	if !exists {
		return confirmedOrders
	}

	confirmedOrders.HallUp = findConfirmedOrdersInArray(orders.HallUpOrders)
	confirmedOrders.HallDown = findConfirmedOrdersInArray(orders.HallDownOrders)
	if cabOrders, exists := orders.CabOrders[ownId]; exists {
		confirmedOrders.Cab = findConfirmedOrdersInArray(cabOrders)
	}

	return confirmedOrders
}
