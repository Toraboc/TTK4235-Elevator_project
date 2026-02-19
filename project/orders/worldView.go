package orders

import (
	. "project/shared"
)
//TODO: Lage no orderhandler og greier med mutex


func CreateWorldView() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	worldView.ConnectedNodes.Add(getOwnId())//circular dependency, will be fixed later

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)

	worldView.Orders.HallUpOrders = CreateOrderList()
	worldView.Orders.HallDownOrders = CreateOrderList()
	worldView.Orders.CabOrders = make(map[NodeId][NumberOfFloors]Order)
	worldView.Orders.CabOrders[getOwnId()] = CreateOrderList()//circular dependency

	return worldView
}

func (wv *WorldView) Copy() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	worldView.ConnectedNodes.Concat(wv.ConnectedNodes)

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)
	for nodeId, elevatorState := range wv.ElevatorStates {
		worldView.ElevatorStates[nodeId] = elevatorState
	}

	worldView.Orders.HallUpOrders = wv.Orders.HallUpOrders.Copy()
	worldView.Orders.HallDownOrders = wv.Orders.HallDownOrders.Copy()

	worldView.Orders.CabOrders = make(map[NodeId][NumberOfFloors]Order)
	for nodeId, cabOrders := range wv.Orders.CabOrders {
		var copiedCabOrders [NumberOfFloors]Order
		for i := 0; i < NumberOfFloors; i++ {
			copiedCabOrders[i] = cabOrders[i].Copy()
		}
		worldView.Orders.CabOrders[nodeId] = copiedCabOrders
	}

	return worldView
}

// This will only sync the orders and elevatorStates
func (worldView *WorldView) Merge(syncMsg SyncMessage) {
	//Oppdaterer staten til heisen m. syncmelding
	worldView.ElevatorStates[syncMsg.Id] = syncMsg.MyState
	//TODO: Sjekk om en caborderliste eksisterer, hvis ikke lag en tom
	//merger ordrer
	for id := range syncMsg.Orders.CabOrders {
		if _, exists := worldView.Orders.CabOrders[id]; !exists {
			worldView.Orders.CabOrders[id] = CreateOrderList()
		}
	}
	for i := 0; i < NumberOfFloors; i++ {
		worldView.Orders.HallDownOrders[i].Merge(syncMsg.Orders.HallDownOrders[i])
		worldView.Orders.HallUpOrders[i].Merge(syncMsg.Orders.HallUpOrders[i])
		for id := range syncMsg.Orders.CabOrders {
			worldView.Orders.CabOrders[id][i].Merge(syncMsg.Orders.CabOrders[id][i])
		}
	}
	
	// This must also be called if our own elevatorsstate changes
	worldView.hallRequestAssigner()
}

// This function will receive updates from the elevator
func (worldView *WorldView) ChangeElevatorState(state ElevatorState) {

	worldView.ElevatorStates[getOwnId()] = state
	worldView.hallRequestAssigner()
}

//run the script and update assigned requests
func (worldView *WorldView) hallRequestAssigner() {

}