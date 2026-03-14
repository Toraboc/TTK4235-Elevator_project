package orderHandler

import (
	"strings"
	. "project/shared"
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
	var copy Orders

	copy.HallUpOrders = orders.HallUpOrders.Clone()
	copy.HallDownOrders = orders.HallDownOrders.Clone()
	copy.CabOrders = make(map[NodeId]*OrderList)
	for nodeId := range orders.CabOrders {
		copy.CabOrders[nodeId] = orders.CabOrders[nodeId].Clone()
	}

	return &copy
}

func (orders *OrderList) Clone() *OrderList {
	var copy OrderList
	for i := range NumberOfFloors {
		copy[i] = orders[i]
	}
	return &copy
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
