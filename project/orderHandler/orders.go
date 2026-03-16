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
	var cloned Orders

	cloned.HallUpOrders = orders.HallUpOrders.Clone()
	cloned.HallDownOrders = orders.HallDownOrders.Clone()
	cloned.CabOrders = make(map[NodeId]*OrderList)
	for nodeId := range orders.CabOrders {
		cloned.CabOrders[nodeId] = orders.CabOrders[nodeId].Clone()
	}

	return &cloned
}

func (orders *OrderList) Clone() *OrderList {
	var cloned OrderList
	for i := range NumberOfFloors {
		cloned[i] = orders[i]
	}
	return &cloned
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
