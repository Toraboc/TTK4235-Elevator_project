package orders

import (
	. "project/shared"
	"time"
)

func CreateOrder(nodeId NodeId) Order {
	order := make(Order)
	order[nodeId] = NO_ORDER
	return order
}

func (order *Order) Copy() Order {
	copy := make(Order)
		for nodeId, status := range *order {
			copy[nodeId] = status
		}
	return copy
}

func CreateOrderList() [NumberOfFloors]Order {
	var orders [NumberOfFloors]Order
	for i := range NumberOfFloors {
		orders[i] = CreateOrder()
	}
	return orders
}
func (orders *[NumberOfFloors]Order) Copy() [NumberOfFloors]Order {
	var copy [NumberOfFloors]Order
	for i := range NumberOfFloors {
		copy[i] = orders[i].Copy()
	}
	return copy
}

//Merge two orders
func (vorldWiewOrder *Order) Merge(syncOrder Order) {

}


// The datainout here will we figure out later
func NewOrder() {

}

// Return the next target floor
func GetNextTargetFloor(worldView *WorldView) (int, error) {
	hallUpOrders := worldView.AssignedHallUpOrders
	hallDownOrders := worldView.AssignedHallDownOrders
	cabOrders := worldView.AssignedCabOrders

	
}
