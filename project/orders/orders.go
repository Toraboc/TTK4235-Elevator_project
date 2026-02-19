package orders

import (
	. "project/shared"
	"time"
)

var worldView WorldView

// Merge our worldview with the incomming data in some way, not dependent on network
// This will only sync the orders and elevatorStates
func MergeWorldView(SyncMessage SyncMessage) {

	// At last
	// This must also be called if our own elevatorsstate changes
	hallRequestAssigner()
}

func GetWorldView() *WorldView {
	return &worldView
}

func UpdateConnectedNodes(ids []NodeId) {
	worldView.ConnectedNodes = ids
}

func CreateOrder() Order {
	var order Order
	order.LastEvent = COMPLETED
	order.ConfirmedBy = make([]NodeId, 0)
	order.LastUpdate = time.Unix(0, 0)
	return order
}

func CreateOrderList() [NumberOfFloors]Order {
	var orders [NumberOfFloors]Order
	for i := range NumberOfFloors {
		orders[i] = CreateOrder()
	}
	return orders
}

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
	return 1, nil
}
