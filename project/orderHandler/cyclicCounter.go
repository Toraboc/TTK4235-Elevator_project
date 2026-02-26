package orderHandler

import (
	"slices"
	. "project/shared"
)

type OrderStatus int
const (
    NO_ORDER OrderStatus = iota
    UNCONFIRMED
    CONFIRMED
    FINISHED
)

func AllEquals[T comparable](slice []T, values []T) bool {
	for _, item := range slice {
		if !slices.Contains(values, item) {
			return false
		}
	}
	return true
}

func getNextValueFromCyclicCounter(myStatus OrderStatus, otherStatuses []OrderStatus) OrderStatus {
	// TODO: This will not allow the state to be updated two steps, this probably needs to be fixed
	switch myStatus {
	case NO_ORDER:
		if slices.Contains(otherStatuses, UNCONFIRMED) {
			return UNCONFIRMED
		}
	case UNCONFIRMED:
		if AllEquals(otherStatuses, []OrderStatus{UNCONFIRMED, CONFIRMED}) {
			return CONFIRMED
		}
	case CONFIRMED:
		if slices.Contains(otherStatuses, FINISHED) {
			return FINISHED
		}
	case FINISHED:
		if AllEquals(otherStatuses, []OrderStatus{FINISHED, NO_ORDER}) {
			return NO_ORDER
		}
	}

	return myStatus
}

type OrderStatusCombined struct {
	myStatus OrderStatus
	otherStatuses []OrderStatus
}

func getOrderStatuses(
	orders map[NodeId]Orders,
	myId NodeId,
	otherNodes NodeIdSet,
	fieldSelector func(Orders) *OrderList,
) [NumberOfFloors]OrderStatusCombined {
	var result [NumberOfFloors]OrderStatusCombined

	for floor := range NumberOfFloors {
		otherStatuses := make([]OrderStatus, len(otherNodes))

		for nodeId, _ := range otherNodes {
			otherStatuses = append(otherStatuses, fieldSelector(orders[nodeId])[floor])
		}

		result[floor].otherStatuses = otherStatuses
		result[floor].myStatus = fieldSelector(orders[myId])[floor]

	}
	return result
}

func updateCyclicCounter(
	orders map[NodeId]Orders,
	myId NodeId,
	otherNodes NodeIdSet,
	fieldSelector func(Orders) *OrderList,
) {
	currentStatus := getOrderStatuses(orders, myId, otherNodes, fieldSelector)
	for floor := range NumberOfFloors {
		nextStatus := getNextValueFromCyclicCounter(currentStatus[floor].myStatus, currentStatus[floor].otherStatuses)
		if nextStatus != currentStatus[floor].myStatus {
			fieldSelector(orders[myId])[floor] = nextStatus
		}
	}
}
