package shared

type OrderCompleted struct {
	Floor int
	Direction Direction
}

type OrderNew struct {
	Floor int
	Type OrderType
}

type OrderType int

const(
	HALLUP 		OrderType = iota
	HALLDOWN
	CAB
)

func (orderType OrderType) String() string {
	switch orderType {
	case HALLUP:
		return "HALLUP"
	case HALLDOWN:
		return "HALLDOWN"
	case CAB:
		return "CAB"
	default:
		panic("Invalid orderType, could not convert to string")
	}
}
