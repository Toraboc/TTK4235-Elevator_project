package shared

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

type ElevatorState struct {
    behaviour  ElevatorBehaviour
    position  int
    direction Direction
}

func (state ElevatorState) Behaviour() ElevatorBehaviour {
	return state.behaviour
}

func (state ElevatorState) Floor() int {
	return state.position
}

func (state ElevatorState) Direction() Direction {
	return state.direction
}
