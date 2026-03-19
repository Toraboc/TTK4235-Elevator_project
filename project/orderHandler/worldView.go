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

func (worldView *WorldView) newOrder(newOrder NewOrderEvent) {
	myId := GetMyId()
	myOrders := worldView.Orders[myId]
	switch newOrder.OrderType {
	case HALLUP:
		if myOrders.HallUpOrders[newOrder.Floor] == NO_ORDER {
			myOrders.HallUpOrders[newOrder.Floor] = UNCONFIRMED
		}
	case HALLDOWN:
		if myOrders.HallDownOrders[newOrder.Floor] == NO_ORDER {
			myOrders.HallDownOrders[newOrder.Floor] = UNCONFIRMED
		}
	case CAB:
		myCabOrders := myOrders.CabOrders[myId]
		if myCabOrders[newOrder.Floor] == NO_ORDER {
			myCabOrders[newOrder.Floor] = UNCONFIRMED
		}
		myOrders.CabOrders[myId] = myCabOrders
	}
}

func (worldView *WorldView) completedOrder(orderCompleted OrderCompletedEvent) {
	myId := GetMyId()
	myOrders := worldView.Orders[myId]

	myOrders.CabOrders[myId][orderCompleted.Floor] = FINISHED

	var hadHallOrder bool
	switch orderCompleted.Direction {
	case UP:
		hadHallOrder = myOrders.HallUpOrders[orderCompleted.Floor] == CONFIRMED
		if hadHallOrder {
			myOrders.HallUpOrders[orderCompleted.Floor] = FINISHED
		}
	case DOWN:
		hadHallOrder = myOrders.HallDownOrders[orderCompleted.Floor] == CONFIRMED
		if hadHallOrder {
			myOrders.HallDownOrders[orderCompleted.Floor] = FINISHED
		}
	}

	if !hadHallOrder {
		targetFloor, err := getNextTargetFloor(*worldView, myId)
		if err != nil {
			panic(err.Error())
		}

		if targetFloor == orderCompleted.Floor {
			switch orderCompleted.Direction {
			case UP:
				myOrders.HallDownOrders[orderCompleted.Floor] = FINISHED
			case DOWN:
				myOrders.HallUpOrders[orderCompleted.Floor] = FINISHED
			}

		}
	}
}

func (worldView *WorldView) handleStateChange() (int, bool) {
	worldView.updateAllOrderStatuses()
	targetFloor, err := getNextTargetFloor(*worldView, GetMyId())
	if err != nil {
		if !strings.HasPrefix(err.Error(), "Missing") {
			panic(err.Error())
		}

		return worldView.lastTargetFloor, false
	}

	if targetFloor != worldView.lastTargetFloor {
		worldView.lastTargetFloor = targetFloor
		return targetFloor, true
	}

	return targetFloor, false
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

	for nodeId := range worldView.ConnectedNodes {
		getCabOrder := func(orders *Orders) *OrderList {
			return orders.CabOrders[nodeId]
		}
		updateCyclicCounter(worldView.Orders, myId, connectedNodes, getCabOrder)
	}
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
