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
	copyOrder := make(Order)
		for nodeId, status := range *order {
			copyOrder[nodeId] = status
		}
	return copyOrder
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
func (vWOrder *Order) Merge(syncOrder Order) {

}

//run the script and update assigned requests
func hallRequestAssigner() {

}

// This function will receive updates from the elevator
func ElevatorStateChange(state ElevatorState) {

}

// The datainout here will we figure out later
func NewOrder() {

}

// Return the next target floor
func GetNextTargetFloor() (int, error) {
	
}
