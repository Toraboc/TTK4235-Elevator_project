package orderHandler

import (
	"fmt"
	. "project/shared"
	"strings"
)

type Orders struct {
	HallUpOrders   *OrderList
	HallDownOrders *OrderList
	CabOrders      map[NodeId]*OrderList
}

func createOrders(nodeId NodeId) *Orders {
	var orders Orders

	orders.HallUpOrders = &OrderList{}
	orders.HallDownOrders = &OrderList{}
	orders.CabOrders = make(map[NodeId]*OrderList)
	orders.CabOrders[nodeId] = &OrderList{}

	return &orders
}

func (orders *Orders) Clone() *Orders {
	var clone Orders

	clone.HallUpOrders = orders.HallUpOrders.Clone()
	clone.HallDownOrders = orders.HallDownOrders.Clone()
	clone.CabOrders = make(map[NodeId]*OrderList)
	for nodeId := range orders.CabOrders {
		clone.CabOrders[nodeId] = orders.CabOrders[nodeId].Clone()
	}

	return &clone
}

func (orders *OrderList) Clone() *OrderList {
	var clone OrderList
	for i := range NumberOfFloors {
		clone[i] = orders[i]
	}
	return &clone
}

func (orders Orders) String() string {
	var builder strings.Builder
	builder.WriteString("Orders{\n")

	fmt.Fprintf(&builder, "\tHallUpOrders: %v,\n", orders.HallUpOrders)
	fmt.Fprintf(&builder, "\tHallDownOrders: %v,\n", orders.HallDownOrders)

	builder.WriteString("\tCabOrders: {\n")
	for nodeId, orderList := range SortedMap(orders.CabOrders) {
		orderListString := strings.ReplaceAll(orderList.String(), "\n", "\n\t\t")
		fmt.Fprintf(&builder, "\t\t[%v]: %s\n", nodeId, orderListString)
	}
	builder.WriteString("\t}\n")

	builder.WriteString("}")
	return builder.String()
}
