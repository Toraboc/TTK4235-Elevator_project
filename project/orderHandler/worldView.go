package orderHandler

import (
	"fmt"
	. "project/shared"
	"strings"
)

type WorldView struct {
	Orders                 map[NodeId]*Orders
	ConnectedNodes         NodeIdSet
	ElevatorStates         map[NodeId]ElevatorState
	AssignedHallUpOrders   [NumberOfFloors]bool
	AssignedHallDownOrders [NumberOfFloors]bool
	AssignedCabOrders      [NumberOfFloors]bool
	lastTargetFloor        int
}

type SyncView struct {
	NodeId        NodeId
	ElevatorState ElevatorState
	Orders        Orders
}

func newWorldView() WorldView {
	var worldView WorldView

	worldView.ConnectedNodes = make(NodeIdSet)
	myId := GetMyId()
	worldView.ConnectedNodes.Add(myId)

	worldView.ElevatorStates = make(map[NodeId]ElevatorState)
	worldView.Orders = make(map[NodeId]*Orders)
	worldView.Orders[myId] = NewOrders(myId)

	worldView.lastTargetFloor = -1

	return worldView
}

func (worldView *WorldView) clone() WorldView {
	var clone WorldView

	clone.ConnectedNodes = make(NodeIdSet)
	clone.ConnectedNodes.Concat(worldView.ConnectedNodes)

	clone.ElevatorStates = make(map[NodeId]ElevatorState)
	for nodeId, elevatorState := range worldView.ElevatorStates {
		clone.ElevatorStates[nodeId] = elevatorState
	}

	clone.Orders = make(map[NodeId]*Orders)
	for nodeId, orders := range worldView.Orders {
		clone.Orders[nodeId] = orders.Clone()
	}

	clone.AssignedHallUpOrders = worldView.AssignedHallUpOrders
	clone.AssignedHallDownOrders = worldView.AssignedHallDownOrders
	clone.AssignedCabOrders = worldView.AssignedCabOrders

	return clone
}

func (worldView *WorldView) merge(sourceNodeId NodeId, sourceNodeState ElevatorState, sourceOrders Orders) {
	worldView.ElevatorStates[sourceNodeId] = sourceNodeState

	worldView.Orders[sourceNodeId] = sourceOrders.Clone()

	cabOrdersIThinkYouHave := worldView.Orders[GetMyId()].CabOrders
	cabOrdersYouHave := sourceOrders.CabOrders
	for nodeId, cabOrders := range cabOrdersYouHave {
		if _, exists := cabOrdersIThinkYouHave[nodeId]; !exists {
			cabOrdersIThinkYouHave[nodeId] = cabOrders.Clone()
		}
	}
}

func (worldView *WorldView) handleStateChange() (int, bool, error) {
	worldView.updateCyclicCounter()

	worldView.hallRequestAssigner()
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

func (worldView *WorldView) updateCyclicCounter() {
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
	if elevatorState.Direction == DOWN {
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

func (worldView *WorldView) String() string {
	var builder strings.Builder

	builder.WriteString("WorldView{\n")
	builder.WriteString("\tConnectedNodes: ")
	builder.WriteString(worldView.ConnectedNodes.String())
	builder.WriteString(",\n")

	builder.WriteString("\tElevatorStates: {\n")
	for nodeId, elevatorState := range SortedMap(worldView.ElevatorStates) {
		builder.WriteString("\t[" + nodeId.String() + "]: ")
		stateString := strings.ReplaceAll(elevatorState.String(), "\n", "\n\t\t")
		builder.WriteString(stateString)
		builder.WriteString("\n")
	}
	builder.WriteString("\t}\n")

	builder.WriteString("\tOrders: {\n")
	for nodeId, orders := range SortedMap(worldView.Orders) {
		builder.WriteString("\t[" + nodeId.String() + "]: ")
		ordersString := strings.ReplaceAll(orders.String(), "\n", "\n\t\t")
		builder.WriteString(ordersString)
		builder.WriteString("\n")
	}
	builder.WriteString("\t}\n")

	builder.WriteString("}")
	return builder.String()
}
