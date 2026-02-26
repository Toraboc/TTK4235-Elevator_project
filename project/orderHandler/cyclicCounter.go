package orderHandler

import (
	. "project/shared"
	"slices"
)

type OrderStatus int

const (
	NO_ORDER OrderStatus = iota
	UNCONFIRMED
	CONFIRMED
	FINISHED
)

func (orderStatus OrderStatus) String() string {
	switch orderStatus {
	case NO_ORDER:
		return "NO ORDER"
	case UNCONFIRMED:
		return "UNCONFIRMED"
	case CONFIRMED:
		return "CONFIRMED"
	case FINISHED:
		return "FINISHED"
	default:
		panic("Invalid orderStatus, could not convert to string")
	}
}

func AllEquals[T comparable](slice []T, values []T) bool {
	for _, item := range slice {
		if !slices.Contains(values, item) {
			return false
		}
	}
	return true
}

func getNextValueFromCyclicCounter(myStatus OrderStatus, connectedNodes []OrderStatus) OrderStatus {
	// TODO: This will not allow the state to be updated two steps, this probably needs to be fixed
	switch myStatus {
	case NO_ORDER:
		if slices.Contains(connectedNodes, UNCONFIRMED) {
			return UNCONFIRMED
		}
	case UNCONFIRMED:
		if AllEquals(connectedNodes, []OrderStatus{UNCONFIRMED, CONFIRMED}) {
			return CONFIRMED
		}
	case CONFIRMED:
		if slices.Contains(connectedNodes, FINISHED) {
			return FINISHED
		}
	case FINISHED:
		if AllEquals(connectedNodes, []OrderStatus{FINISHED, NO_ORDER}) {
			return NO_ORDER
		}
	}

	return myStatus
}

type OrderStatusCombined struct {
	myStatus      OrderStatus
	otherStatuses []OrderStatus
}

func getOrderStatuses(
	orders map[NodeId]Orders,
	myId NodeId,
	connectedNodes NodeIdSet,
	fieldSelector func(Orders) *OrderList,
) [NumberOfFloors]OrderStatusCombined {
	var result [NumberOfFloors]OrderStatusCombined

	for floor := range NumberOfFloors {
		otherStatuses := make([]OrderStatus, len(connectedNodes))

		i := 0
		for nodeId, _ := range connectedNodes {
			otherStatuses[i] = fieldSelector(orders[nodeId])[floor]
			i++
		}

		result[floor].otherStatuses = otherStatuses
		result[floor].myStatus = fieldSelector(orders[myId])[floor]

	}
	return result
}

func updateCyclicCounter(
	orders map[NodeId]Orders,
	myId NodeId,
	connectedNodes NodeIdSet,
	fieldSelector func(Orders) *OrderList,
) {
	currentStatus := getOrderStatuses(orders, myId, connectedNodes, fieldSelector)
	for floor := range NumberOfFloors {
		nextStatus := getNextValueFromCyclicCounter(currentStatus[floor].myStatus, currentStatus[floor].otherStatuses)
		if nextStatus != currentStatus[floor].myStatus {
			fieldSelector(orders[myId])[floor] = nextStatus
		}
	}
}
