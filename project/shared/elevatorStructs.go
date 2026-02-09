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
)

func (behaviour ElevatorBehaviour) CanBeAssignedOrders() bool {
	return behaviour == IDLE || behaviour == MOVING || behaviour == PASSENGER_TRANSFER
}

func (behaviour ElevatorBehaviour) CanMove() bool {
	return behaviour == MOVING || behaviour == FAULTY_MOTOR
}

type ElevatorState struct {
    bahviour  ElevatorBehaviour
    position  int
    direction Direction
}
