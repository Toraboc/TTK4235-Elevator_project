package orders

import (
	"time"
	. "project/shared"
)
//TODO: Lage no orderhandler og greier med mutex


func CreateWorldView() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	worldView.ConnectedNodes.Add(GetMyId())

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)

	worldView.Orders.HallUpOrders = CreateOrderList()
	worldView.Orders.HallDownOrders = CreateOrderList()
	worldView.Orders.CabOrders = make(map[NodeId]OrderList)
	worldView.Orders.CabOrders[GetMyId()] = CreateOrderList()

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

	worldView.Orders.CabOrders = make(map[NodeId]OrderList)
	for nodeId, cabOrders := range wv.Orders.CabOrders {
		var copiedCabOrders OrderList
		for i := 0; i < NumberOfFloors; i++ {
			copiedCabOrders[i] = cabOrders[i].Copy()
		}
		worldView.Orders.CabOrders[nodeId] = copiedCabOrders
	}

	return worldView
}


//NB DENNE ER NÅ TO STEDER; MEN STRUKTUR MÅ ENDRES SLIK AT DENNE BLIR TILGJENGELIG UTEN SIRKULÆRE PAKKEGREIER
type SyncMessage struct {
	Id         NodeId
	Orders     Orders
	MyState    ElevatorState
	KnownNodes []NodeId
	SendTime   time.Time
}


// This will only sync the orders and elevatorStates
func (worldView *WorldView) Merge(syncMsg *SyncMessage) {
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
		worldView.Orders.HallDownOrders[i].Merge(syncMsg.Orders.HallDownOrders[i],syncMsg.Id)
		worldView.Orders.HallUpOrders[i].Merge(syncMsg.Orders.HallUpOrders[i],syncMsg.Id)
		for id := range syncMsg.Orders.CabOrders {
			cabOrderList := worldView.Orders.CabOrders[id]
			cabOrder := cabOrderList[i]
			cabOrder.Merge(syncMsg.Orders.CabOrders[id][i], syncMsg.Id)
			cabOrderList[i] = cabOrder
			worldView.Orders.CabOrders[id] = cabOrderList
		}
	}
	//TODO: update cyclic counter
	
	
	// This must also be called if our own elevatorsstate changes
	worldView.hallRequestAssigner()
}

// This function will receive updates from the elevator
func (worldView *WorldView) ChangeElevatorState(state ElevatorState) {

	worldView.ElevatorStates[GetMyId()] = state
	worldView.hallRequestAssigner()
}

