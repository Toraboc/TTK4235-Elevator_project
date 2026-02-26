package orderHandler

import (
	"fmt"
	. "project/shared"
)


type OrderList [NumberOfFloors]OrderStatus

type Orders struct {
	HallUpOrders   OrderList
	HallDownOrders OrderList
	CabOrders      map[NodeId]OrderList
}

func newOrders(nodeId NodeId) Orders {
	var orders Orders

	orders.CabOrders = make(map[NodeId]OrderList)
	orders.CabOrders[nodeId] = OrderList{}

	return orders
}

func (orders *Orders) clone() Orders {
	var copy Orders

	copy.HallUpOrders = orders.HallUpOrders.clone()
	copy.HallDownOrders = orders.HallDownOrders.clone()
	copy.CabOrders = make(map[NodeId]OrderList)
	for nodeId := range orders.CabOrders {
		copy.CabOrders[nodeId] = orders.CabOrders[nodeId].clone()
	}

	return copy
}

func (orders OrderList) clone() OrderList {
	var copy OrderList
	for i := range NumberOfFloors {
		copy[i] = orders[i]
	}
	return copy
}


func (worldView *WorldView) getNextTargetFloor() (int, error) {          
	//Feilsøkingsgreier som kan fjernes etterhvert
	if worldView == nil {
		return -1, fmt.Errorf("worldView is nil")
	}

	myId := GetMyId()

	elevatorState, exists := worldView.ElevatorStates[myId]
	if !exists {
		return -1, fmt.Errorf("missing elevator elevatorState for own node")
	}

	floor := elevatorState.Floor
	if floor < 0 || floor >= NumberOfFloors {
		return -1, fmt.Errorf("invalid current floor: %d", floor)
	}

	myOrders := worldView.Orders[myId]

	if elevatorState.Direction == UP {
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if myOrders.CabOrders[myId][floor] == CONFIRMED || myOrders.HallUpOrders[floor] == CONFIRMED {
				return floor, nil
			}
		}

		for i := floor + 1; i < NumberOfFloors; i++ {
			if myOrders.CabOrders[myId][i] == CONFIRMED || myOrders.HallUpOrders[i] == CONFIRMED {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= 0; i-- {
			if myOrders.CabOrders[myId][i] == CONFIRMED || myOrders.HallDownOrders[i] == CONFIRMED {
				return i, nil
			}
		}
		for i := 0; i <= floor; i++ {
			if myOrders.CabOrders[myId][i] == CONFIRMED || myOrders.HallUpOrders[i] == CONFIRMED {
				return i, nil
			}
		}
	}
	//Denne er vel strengt talt ikke nødvendig, men grei for ryddighetens skyld
	if elevatorState.Direction == DOWN{
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if myOrders.CabOrders[myId][floor] == CONFIRMED || myOrders.HallDownOrders[floor] == CONFIRMED {
				return floor, nil
			}
		}

		for i := floor - 1; i >= 0; i-- {
			if myOrders.CabOrders[myId][i] == CONFIRMED || myOrders.HallDownOrders[i] == CONFIRMED {
				return i, nil
			}
		}
		for i := range NumberOfFloors {
			if myOrders.CabOrders[myId][i] == CONFIRMED || myOrders.HallUpOrders[i] == CONFIRMED {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= floor; i-- {
			if myOrders.CabOrders[myId][i] == CONFIRMED || myOrders.HallDownOrders[i] == CONFIRMED {
				return i, nil
			}
		}
	}
	return -1, nil
}
