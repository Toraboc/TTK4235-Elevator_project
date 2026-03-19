package orderHandler

import (
	. "project/shared"
	"slices"
)

type OrderStatusCombined struct {
	myStatus      OrderStatus
	otherStatuses []OrderStatus
}

func updateCyclicCounter(
	orders map[NodeId]*Orders,
	myId NodeId,
	connectedNodes NodeIdSet,
	fieldSelector func(*Orders) *OrderList,
) {
	currentStatus := getOrderStatuses(orders, myId, connectedNodes, fieldSelector)
	for floor := range NumberOfFloors {
		nextStatus := cyclicCounterNextValue(currentStatus[floor].myStatus, currentStatus[floor].otherStatuses)
		if nextStatus != currentStatus[floor].myStatus {
			orderList := fieldSelector(orders[myId])
			orderList[floor] = nextStatus
		}
	}
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

		for nodeId := range connectedNodes {
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
		if orderList != nil {
			result[floor].myStatus = orderList[floor]
		}

	}
	return result
}

func cyclicCounterNextValue(myStatus OrderStatus, connectedNodeStatuses []OrderStatus) OrderStatus {
	switch myStatus {
	case NO_ORDER:
		if slices.Contains(connectedNodeStatuses, CONFIRMED) || (allEquals(connectedNodeStatuses, []OrderStatus{UNCONFIRMED}) && len(connectedNodeStatuses) > 0) {
			return CONFIRMED
		}
		if slices.Contains(connectedNodeStatuses, UNCONFIRMED) {
			return UNCONFIRMED
		}
	case UNCONFIRMED:
		if allEquals(connectedNodeStatuses, []OrderStatus{UNCONFIRMED, CONFIRMED}) || slices.Contains(connectedNodeStatuses, CONFIRMED) {
			return CONFIRMED
		}
		if slices.Contains(connectedNodeStatuses, FINISHED) {
			return FINISHED
		}
	case CONFIRMED:
		if allEquals(connectedNodeStatuses, []OrderStatus{FINISHED}) && len(connectedNodeStatuses) > 0 {
			return NO_ORDER
		}
		if slices.Contains(connectedNodeStatuses, FINISHED) {
			return FINISHED
		}
	case FINISHED:
		if allEquals(connectedNodeStatuses, []OrderStatus{FINISHED, NO_ORDER}) {
			return NO_ORDER
		}
	}

	return myStatus
}

func allEquals[T comparable](slice []T, values []T) bool {
	for _, item := range slice {
		if !slices.Contains(values, item) {
			return false
		}
	}
	return true
}
