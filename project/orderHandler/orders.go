package orderHandler

import (
	. "project/shared"
	"strings"
)

type Orders struct {
	HallUpOrders   *OrderList
	HallDownOrders *OrderList
	CabOrders      map[NodeId]*OrderList
}

func NewOrders(nodeId NodeId) *Orders {
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
	builder.WriteString("\tHallUpOrders: ")
	builder.WriteString(orders.HallUpOrders.String())
	builder.WriteString(",\n")

	builder.WriteString("\tHallDownOrders: ")
	builder.WriteString(orders.HallDownOrders.String())
	builder.WriteString(",\n")

	builder.WriteString("\tCabOrders: {\n")
	for nodeId, orderList := range SortedMap(orders.CabOrders) {
		builder.WriteString("\t[" + nodeId.String() + "]: ")
		orderListString := strings.ReplaceAll(orderList.String(), "\n", "\n\t\t")
		builder.WriteString(orderListString)
		builder.WriteString("\n")
	}
	builder.WriteString("\t}\n")

	builder.WriteString("}")
	return builder.String()
}
