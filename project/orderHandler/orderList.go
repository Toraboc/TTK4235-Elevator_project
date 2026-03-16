package orderHandler

import (
	. "project/shared"
	"strings"
)

type OrderList [NumberOfFloors]OrderStatus

func (orderList OrderList) String() string {
	var builder strings.Builder

	builder.WriteString("[")

	for i := range NumberOfFloors {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(orderList[i].String())
	}

	builder.WriteString("]")
	return builder.String()
}
