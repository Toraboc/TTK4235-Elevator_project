package orderHandler

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

