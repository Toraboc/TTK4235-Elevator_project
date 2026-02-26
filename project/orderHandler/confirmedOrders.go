package orderHandler

import (
	. "project/shared"
)

type ConfirmedOrders struct {
	HallUp   [NumberOfFloors]bool
	HallDown [NumberOfFloors]bool
	Cab      [NumberOfFloors]bool
}

func findConfirmedOrdersInArray(orders *OrderList) [NumberOfFloors]bool {
	var confirmed [NumberOfFloors]bool

	for floor := range NumberOfFloors {
		confirmed[floor] = orders[floor] == CONFIRMED
	}
	return confirmed
}

func (worldView *WorldView) getConfirmedOrders() ConfirmedOrders {
	var confirmedOrders ConfirmedOrders
	myId := GetMyId()

	orders, exists := worldView.Orders[myId]
	if !exists {
		panic("No orders, invalid state")
	}

	confirmedOrders.HallUp = findConfirmedOrdersInArray(orders.HallUpOrders)
	confirmedOrders.HallDown = findConfirmedOrdersInArray(orders.HallDownOrders)
	if cabOrders, exists := orders.CabOrders[myId]; exists {
		confirmedOrders.Cab = findConfirmedOrdersInArray(cabOrders)
	}

	return confirmedOrders
}
