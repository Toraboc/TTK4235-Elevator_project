package orderHandler

import (
	. "project/shared"
	"slices"
)

type OrderStatusCombined struct {
	myStatus      OrderStatus
	otherStatuses []OrderStatus
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
	switch myStatus {
	case NO_ORDER:
		if slices.Contains(connectedNodes, CONFIRMED) {
			return CONFIRMED
		}
		if slices.Contains(connectedNodes, UNCONFIRMED) {
			return UNCONFIRMED
		}
	case UNCONFIRMED:
		if AllEquals(connectedNodes, []OrderStatus{UNCONFIRMED, CONFIRMED}) || slices.Contains(connectedNodes, CONFIRMED) {
			return CONFIRMED
		}
	case CONFIRMED:
		if slices.Contains(connectedNodes, FINISHED) {
			return FINISHED
		}
	case FINISHED:
		if slices.Contains(connectedNodes, UNCONFIRMED) {
			return UNCONFIRMED
		}
		if AllEquals(connectedNodes, []OrderStatus{FINISHED, NO_ORDER}) {
			return NO_ORDER
		}
	}

	return myStatus
}

func getOrderStatuses(
	orders map[NodeId]*Orders,
	myId NodeId,
	connectedNodes NodeIdSet,
	fieldSelector func(*Orders) *OrderList,
) [NumberOfFloors]OrderStatusCombined {
	var result [NumberOfFloors]OrderStatusCombined

	for floor := range NumberOfFloors {
		otherStatuses := make([]OrderStatus, 0)

		for nodeId, _ := range connectedNodes {
			nodeOrders, exists := orders[nodeId]
			if !exists {
				continue
			}
			orderList := fieldSelector(nodeOrders)
			if orderList == nil {
				continue
			}
			otherStatuses = append(otherStatuses, orderList[floor])
		}

		result[floor].otherStatuses = otherStatuses

		orderList := fieldSelector(orders[myId])
		// TODO: This is fixing the symptom, not the cause
		if orderList != nil {
			result[floor].myStatus = orderList[floor]
		}

	}
	return result
}

func updateCyclicCounter(
	orders map[NodeId]*Orders,
	myId NodeId,
	connectedNodes NodeIdSet,
	fieldSelector func(*Orders) *OrderList,
) {
	currentStatus := getOrderStatuses(orders, myId, connectedNodes, fieldSelector)
	for floor := range NumberOfFloors {
		nextStatus := getNextValueFromCyclicCounter(currentStatus[floor].myStatus, currentStatus[floor].otherStatuses)
		if nextStatus != currentStatus[floor].myStatus {
			orderList := fieldSelector(orders[myId])
			orderList[floor] = nextStatus
		}
	}
}
