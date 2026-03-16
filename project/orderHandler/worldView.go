package orderHandler

import (
	"fmt"
	"maps"
	. "project/shared"
	"strings"
)

type WorldView struct {
	Orders          map[NodeId]*Orders
	ConnectedNodes  NodeIdSet
	ElevatorStates  map[NodeId]ElevatorState
	lastTargetFloor int
}

func newWorldView() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	myId := GetMyId()
	worldView.ConnectedNodes.Add(myId)

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)
	worldView.Orders = make(map[NodeId]*Orders)
	worldView.Orders[myId] = createOrders(myId)

	worldView.lastTargetFloor = -1

	return worldView
}

func (worldView *WorldView) clone() WorldView {
	var clone WorldView

	clone.ConnectedNodes = make(NodeIdSet)
	clone.ConnectedNodes.Concat(worldView.ConnectedNodes)

	clone.ElevatorStates = maps.Clone(worldView.ElevatorStates)

	clone.Orders = make(map[NodeId]*Orders)
	for nodeId, orders := range worldView.Orders {
		clone.Orders[nodeId] = orders.Clone()
	}

	return clone
}

func (worldView *WorldView) merge(sourceNodeId NodeId, sourceNodeState ElevatorState, sourceOrders Orders) {
	worldView.ElevatorStates[sourceNodeId] = sourceNodeState

	worldView.Orders[sourceNodeId] = sourceOrders.Clone()

	existingCabOrders := worldView.Orders[GetMyId()].CabOrders
	incomingCabOrders := sourceOrders.CabOrders
	for nodeId, cabOrders := range incomingCabOrders {
		if _, exists := existingCabOrders[nodeId]; !exists {
			existingCabOrders[nodeId] = cabOrders.Clone()
		}
	}
}

func (worldView *WorldView) handleStateChange() (int, bool, error) {
	worldView.updateAllOrderStatuses()
	targetFloor, err := worldView.getNextTargetFloor()
	if err != nil {
		fmt.Println(err.Error())
		return worldView.lastTargetFloor, false, err
	}

	if targetFloor != worldView.lastTargetFloor {
		worldView.lastTargetFloor = targetFloor
		return targetFloor, true, nil
	}

	return targetFloor, false, nil
}

func (worldView *WorldView) updateAllOrderStatuses() {
	myId := GetMyId()
	connectedNodes := worldView.ConnectedNodes.Clone()
	connectedNodes.Remove(myId)

	getHallDown := func(orders *Orders) *OrderList {
		return orders.HallDownOrders
	}
	updateCyclicCounter(worldView.Orders, myId, connectedNodes, getHallDown)

	getHallUp := func(orders *Orders) *OrderList {
		return orders.HallUpOrders
	}
	updateCyclicCounter(worldView.Orders, myId, connectedNodes, getHallUp)

	for nodeId := range connectedNodes {
		getCabOrder := func(orders *Orders) *OrderList {
			return orders.CabOrders[nodeId]
		}
		updateCyclicCounter(worldView.Orders, myId, connectedNodes, getCabOrder)
	}

	getMyCab := func(orders *Orders) *OrderList {
		return orders.CabOrders[myId]
	}
	updateCyclicCounter(worldView.Orders, myId, connectedNodes, getMyCab)
}

func (worldView *WorldView) getNextTargetFloor() (int, error) {
	myId := GetMyId()

	elevatorState, exists := worldView.ElevatorStates[myId]
	if !exists {
		return -1, fmt.Errorf("missing elevator elevatorState for own node")
	}

	floor := elevatorState.Floor
	if floor < 0 || floor >= NumberOfFloors {
		return -1, fmt.Errorf("invalid current floor: %d", floor)
	}

	orders := hallRequestAssigner(worldView, myId)

	if elevatorState.Direction == UP {
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if orders.Cab[floor] || orders.HallUp[floor] {
				return floor, nil
			}
		}

		for i := floor + 1; i < NumberOfFloors; i++ {
			if orders.Cab[i] || orders.HallUp[i] {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= 0; i-- {
			if orders.Cab[i] || orders.HallDown[i] {
				return i, nil
			}
		}
		for i := 0; i <= floor; i++ {
			if orders.Cab[i] || orders.HallUp[i] {
				return i, nil
			}
		}
	}
	if elevatorState.Direction == DOWN {
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if orders.Cab[floor] || orders.HallDown[floor] {
				return floor, nil
			}
		}

		for i := floor - 1; i >= 0; i-- {
			if orders.Cab[i] || orders.HallDown[i] {
				return i, nil
			}
		}
		for i := range NumberOfFloors {
			if orders.Cab[i] || orders.HallUp[i] {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= floor; i-- {
			if orders.Cab[i] || orders.HallDown[i] {
				return i, nil
			}
		}
	}
	return -1, nil
}

func (worldView *WorldView) String() string {
	var builder strings.Builder
	builder.WriteString("WorldView{\n")
	fmt.Fprintf(&builder, "\tConnectedNodes: %v,\n", worldView.ConnectedNodes)

	builder.WriteString("\tElevatorStates: {\n")
	for nodeId, elevatorState := range SortedMap(worldView.ElevatorStates) {
		stateString := strings.ReplaceAll(elevatorState.String(), "\n", "\n\t\t")
		fmt.Fprintf(&builder, "\t\t[%v]: %s\n", nodeId, stateString)
	}
	builder.WriteString("\t}\n")

	builder.WriteString("\tOrders: {\n")
	for nodeId, orders := range SortedMap(worldView.Orders) {
		ordersString := strings.ReplaceAll(orders.String(), "\n", "\n\t\t")
		fmt.Fprintf(&builder, "\t\t[%v]: %s\n", nodeId, ordersString)
	}
	builder.WriteString("\t}\n")

	builder.WriteString("}")
	return builder.String()
}
