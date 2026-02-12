package orders

import (
	. "project/shared"
)

var worldView WorldView
//todo: Initialize worldview for the first time
//		set all orders to completed with timestamp 0
func initializeWorldView(worldView WorldView){

}

// This will only sync the orders and elevatorStates
func MergeWorldView(syncMsg SyncMessage) {
	//Oppdaterer staten til heisen m. syncmelding
	worldView.ElevatorStates[syncMsg.Id] = syncMsg.MyState

	//merger ordrer
	for i := 0; i < NumberOfFloors; i++ {
		if i != 0{
			worldView.Orders.HallDownOrders[i].MergeOrder(syncMsg.Orders.HallDownOrders[i])
		}
		if i != (NumberOfFloors -1){
			worldView.Orders.HallUpOrders[i].MergeOrder(syncMsg.Orders.HallUpOrders[i])
		}
		for id := range syncMsg.Orders.CabOrders {
			worldView.Orders.CabOrders[id][i].MergeOrder(syncMsg.Orders.CabOrders[id][i])
		}
	}
	
	// This must also be called if our own elevatorsstate changes
	hallRequestAssigner()
}
func (vWOrder *Order) MergeOrder(syncOrder Order) {
	if vWOrder.LastEvent == syncOrder.LastEvent {
		// select newest timestamp when minor disagreements occur
		if syncOrder.LastUpdate.After(vWOrder.LastUpdate) {
			vWOrder.LastUpdate = syncOrder.LastUpdate 
		}
		vWOrder.ConfirmedBy.Concat(syncOrder.ConfirmedBy)
	} else {
		if syncOrder.LastUpdate.After(vWOrder.LastUpdate) {
			vWOrder.LastUpdate = syncOrder.LastUpdate 
			vWOrder.ConfirmedBy = syncOrder.ConfirmedBy
			vWOrder.ConfirmedBy.Add(getOwnId)
		}
	}
}

func GetWorldview() WorldView {
	return worldView
}
//run the script and update assigned requests
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
