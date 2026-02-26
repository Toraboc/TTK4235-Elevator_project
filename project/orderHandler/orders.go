package orderHandler

import (
	"fmt"
	. "project/shared"
	"strings"
)


type OrderList [NumberOfFloors]OrderStatus

func (orderList OrderList) String() string {
	var builder strings.Builder

	builder.WriteString("[")

	for i := range NumberOfFloors {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(orderList[i].String())
	}

	builder.WriteString("]")
	return builder.String()
}

type OrderType int

const(
	HALLUP 		OrderType = iota
	HALLDOWN
	CAB
)


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

	hallUpOrders := worldView.AssignedHallUpOrders
	hallDownOrders := worldView.AssignedHallDownOrders
	cabOrders := worldView.AssignedCabOrders

	if elevatorState.Direction == UP {
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if cabOrders[floor] || hallUpOrders[floor] {
				return floor, nil
			}
		}

		for i := floor + 1; i < NumberOfFloors; i++ {
			if cabOrders[i] || hallUpOrders[i] {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= 0; i-- {
			if cabOrders[i] || hallDownOrders[i] {
				return i, nil
			}
		}
		for i := 0; i <= floor; i++ {
			if cabOrders[i] || hallUpOrders[i] {
				return i, nil
			}
		}
	}
	//Denne er vel strengt talt ikke nødvendig, men grei for ryddighetens skyld
	if elevatorState.Direction == DOWN{
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if cabOrders[floor] || hallDownOrders[floor] {
				return floor, nil
			}
		}

		for i := floor - 1; i >= 0; i-- {
			if cabOrders[i] || hallDownOrders[i] {
				return i, nil
			}
		}
		for i := range NumberOfFloors {
			if cabOrders[i] || hallUpOrders[i] {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= floor; i-- {
			if cabOrders[i] || hallDownOrders[i] {
				return i, nil
			}
		}
	}
	return -1, nil
}

func (orders Orders) String() string {
	var builder strings.Builder

	builder.WriteString("Orders{\n")
	builder.WriteString("\tHallUpOrders: ")
	builder.WriteString(orders.HallUpOrders.String())
	builder.WriteString(",\n")

	builder.WriteString("\tHallDownOrders: ")
	builder.WriteString(orders.HallDownOrders.String())
	builder.WriteString(",\n")

	builder.WriteString("\tCabOrders: {\n")
	for nodeId, orderList := range orders.CabOrders {
		builder.WriteString("\t[" + nodeId.String() + "]: ")
		orderListString := strings.ReplaceAll(orderList.String(), "\n", "\n\t\t")
		builder.WriteString(orderListString)
		builder.WriteString("\n")
	}
	builder.WriteString("\t}\n")

	builder.WriteString("}")
	return builder.String()
}
