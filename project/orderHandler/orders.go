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

func NewOrders(nodeId NodeId) Orders {
	var orders Orders

	orders.HallUpOrders = NewOrderList()
	orders.HallDownOrders = NewOrderList()
	orders.CabOrders = make(map[NodeId]OrderList)
	orders.CabOrders[nodeId] = NewOrderList()

	return orders
}
func (orders *Orders) Clone() Orders {
	var copy Orders

	copy.HallUpOrders = orders.HallUpOrders.Clone()
	copy.HallDownOrders = orders.HallDownOrders.Clone()
	copy.CabOrders = make(map[NodeId]OrderList)
	for nodeId := range orders.CabOrders {
		copy.CabOrders[nodeId] = orders.CabOrders[nodeId].Clone()
	}

	return copy
}
 
func NewOrderList() OrderList {
	var orders OrderList
	for i := range NumberOfFloors {
		orders[i] = NO_ORDER
	}
	return orders
}



func (orders OrderList) Clone() OrderList {
	var copy OrderList
	for i := range NumberOfFloors {
		copy[i] = orders[i]
	}
	return copy
}


func (worldView *WorldView) UpdateCounter() {

	
}

// The datainout here will we figure out later
func NewOrder() {

}

func GetNextTargetFloor(worldView *WorldView) (int, error) {          
	//Feilsøkingsgreier som kan fjernes etterhvert
	if worldView == nil {
		return -1, fmt.Errorf("worldView is nil")
	}

	elevatorState, exists := worldView.ElevatorStates[GetMyId()]
	if !exists {
		return -1, fmt.Errorf("missing elevator elevatorState for own node")
	}

	floor := elevatorState.Floor()
	if floor < 0 || floor >= NumberOfFloors {
		return -1, fmt.Errorf("invalid current floor: %d", floor)
	}

	

	if elevatorState.Direction() == UP {
		if elevatorState.Behaviour() == PASSENGER_TRANSFER || elevatorState.Behaviour() == IDLE {
			if hasUpRequestAtFloor(worldView, floor) {
				return floor, nil
			}
			if elevatorState.Behaviour() == IDLE {
				if hasDownRequestAtFloor(worldView, floor) {
					return floor, nil
				}
			}
		}
		if target := nearestUpRequestAbove(worldView, floor); target != -1 {
			return target, nil
		}
		if target := nearestDownRequestAbove(worldView, floor); target != -1 {
			return target, nil
		}
		if hasDownRequestAtFloor(worldView, floor) {
			return floor, nil
		}
		if target := nearestDownRequestBelow(worldView, floor); target != -1 {
			return target, nil
		}
		if target := nearestUpRequestBelow(worldView, floor); target != -1 {
			return target, nil
		}
		if hasUpRequestAtFloor(worldView, floor) {
			return floor, nil
		}
		return -1, nil
	}
	//Denne er vel strengt talt ikke nødvendig, men grei for ryddighetens skyld
	if elevatorState.Direction() == DOWN{
		if elevatorState.Behaviour() == PASSENGER_TRANSFER || elevatorState.Behaviour() == IDLE {
			if hasDownRequestAtFloor(worldView, floor) {
				return floor, nil
			}
			if elevatorState.Behaviour() == IDLE {
				if hasUpRequestAtFloor(worldView, floor) {
					return floor, nil
				}
			}
		}
		if target := nearestDownRequestBelow(worldView, floor); target != -1 {
			return target, nil
		}
		if target := nearestUpRequestBelow(worldView, floor); target != -1 {
			return target, nil
		}
		if hasUpRequestAtFloor(worldView, floor) {
			return floor, nil
		}
		if target := nearestUpRequestAbove(worldView, floor); target != -1 {
			return target, nil
		}
		if target := nearestDownRequestAbove(worldView, floor); target != -1 {
			return target, nil
		}
		if hasUpRequestAtFloor(worldView, floor) {
			return floor, nil
		}
		return -1, nil
	}
	return -1, nil
}


func hasUpRequestAtFloor(worldView *WorldView, floor int) bool {
	return worldView.AssignedHallUpOrders[floor] || worldView.AssignedCabOrders[floor]
}

func hasDownRequestAtFloor(worldView *WorldView, floor int) bool {
	return worldView.AssignedHallDownOrders[floor] || worldView.AssignedCabOrders[floor]
}
func hasRequestAtFloor(worldView *WorldView, floor int) bool {
	return hasDownRequestAtFloor(worldView, floor) || hasUpRequestAtFloor(worldView, floor)
}

func nearestUpRequestAbove(worldView *WorldView, floor int) int {
	if floor == NumberOfFloors -1{//Kunne vel vært -2 her
		return -1
	}
	for f := floor + 1; f < NumberOfFloors; f++ {
		if hasUpRequestAtFloor(worldView, f) {
			return f
		}
	}
	return -1
}
func nearestDownRequestAbove(worldView *WorldView, floor int) int {
	if floor == NumberOfFloors -1{
		return -1
	}
	for f := floor + 1; f < NumberOfFloors; f++ {
		if hasDownRequestAtFloor(worldView, f) {
			return f
		}
	}
	return -1
}

func nearestDownRequestBelow(worldView *WorldView, floor int) int {
	if floor == 0{//Kunne vært 1 istedet for 0 her
		return -1
	}
	for f := floor - 1; f >= 0; f-- {
		if hasDownRequestAtFloor(worldView, f) {
			return f
		}
	}
	return -1
}

func nearestUpRequestBelow(worldView *WorldView, floor int) int {
	if floor == 0{
		return -1
	}
	for f := floor - 1; f >= 0; f-- {
		if hasUpRequestAtFloor(worldView, f){
			return f
		}
	}
	return -1
}
