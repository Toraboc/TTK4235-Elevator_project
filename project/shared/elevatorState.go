package shared

type ElevatorState struct {
	Behaviour ElevatorBehaviour
	Floor  int
	Direction Direction
}

type Direction int

const (
	UP Direction = iota
	DOWN
)

type ElevatorBehaviour int

const (
	IDLE ElevatorBehaviour = iota
	MOVING
	PASSENGER_TRANSFER
	DOOR_OBSTRUCTED
	FAULTY_MOTOR
	DISCONNECTED
)


func (behaviour ElevatorBehaviour) CanBeAssignedOrders() bool {
	return behaviour == IDLE || behaviour == MOVING || behaviour == PASSENGER_TRANSFER
}

func (behaviour ElevatorBehaviour) CanMove() bool {
	return behaviour == MOVING || behaviour == FAULTY_MOTOR
}

func (behaviour ElevatorBehaviour) String() string {
	switch behaviour {
	case IDLE:
		return "IDLE"
	case MOVING:
		return "MOVING"
	case PASSENGER_TRANSFER:
		return "PASSENGER_TRANSFER"
	case FAULTY_MOTOR:
		return "FAULTY_MOTOR"
	case DOOR_OBSTRUCTED:
		return "DOOR_OBSTRUCTED"
	case DISCONNECTED:
		return "DISCONNECTED"
	default:
		panic("Undefined EleavtorBehaviour")
	}
}
