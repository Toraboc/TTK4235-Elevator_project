package orders

import (
	. "project/shared"
)

var worldView WorldView
//todo: Initialize worldview for the first time
//		set all orders to completed with timestamp 0
func initializeWorldView(worldView WorldView){

}

	worldView.ElevatorStates[SyncMessage.Id] = SyncMessage.MyState
// This will only sync the orders and elevatorStates
func MergeWorldView(syncMsg SyncMessage) {
	//merger først elevatorStates
	worldView.ElevatorStates[syncMsg.Id] = syncMsg.MyState

	//merger så hallordrer
	//TODO: Bruk MergeOrder her
	for i := 0; i < NumberOfFloors; i++ {
		if i != 0{
			worldView.Orders.HallDownOrders
		}

	}
	
	//slettes snart, structene for oversikt
	/*
	type WorldView struct {
    	Orders Orders
		ConnectedNodes []nodeId
    	ElevatorStates map[NodeId]ElevatorState
    	AssignedHallUpOrders [NumberOfFloors]bool
    	AssignedHallDownOrders [NumberOfFloors]bool
    	AssignedCabOrders [NumberOfFloors]bool
	}
	type SyncMessage struct {
		Id         NodeId
		Orders     Orders
		MyState    ElevatorState
		KnownNodes []NodeId
	}
	*/
	// At last
	// This must also be called if our own elevatorsstate changes
	hallRequestAssigner()
}
func (order1 Order) MergeOrder(order2 Order) Order {
	// TODO: define merge rules for orders
	if order1.LastEvent==order2.LastEvent
	
	return order1
}
/*    Slettes straks
type Order struct {
    LastEvent OrderStatus // skal dette vere OrderStatus?
    LastUpdate time.Time
    ConfirmedBy []NodeId
}*/
func GetWorldview() WorldView {
	return WorldView{}
}

func hallRequestAssigner() {

}

// This function will receive updates from the elevator
func ElevatorStateChange(state ElevatorState) {

}

// The datainout here will we figure out later
func newOrder() {

}

// Return the next target floor
func GetNextTargetFloor() (int, error) {
	return 1, nil
}
