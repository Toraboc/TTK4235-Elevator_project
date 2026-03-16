package shared

import (
	"fmt"
	"strings"
)

type Direction int

const (
	UP Direction = iota
	DOWN
)

func (direction Direction) String() string {
	switch direction {
	case UP:
		return "UP"
	case DOWN:
		return "DOWN"
	default:
		panic("Invalid direction, cannot crete string")
	}
}

type ElevatorBehaviour int

const (
	IDLE ElevatorBehaviour = iota
	MOVING
	PASSENGER_TRANSFER
	DOOR_OBSTRUCTED
	MOTOR_FAILURE
)

func (behaviour ElevatorBehaviour) CanBeAssignedOrders() bool {
	return behaviour == IDLE || behaviour == MOVING || behaviour == PASSENGER_TRANSFER
}

func (behaviour ElevatorBehaviour) Moving() bool {
	return behaviour == MOVING || behaviour == MOTOR_FAILURE
}

func (behaviour ElevatorBehaviour) String() string {
	switch behaviour {
	case IDLE:
		return "IDLE"
	case MOVING:
		return "MOVING"
	case PASSENGER_TRANSFER:
		return "PASSENGER_TRANSFER"
	case MOTOR_FAILURE:
		return "MOTOR_FAILURE"
	case DOOR_OBSTRUCTED:
		return "DOOR_OBSTRUCTED"
	default:
		panic("Undefined ElevatorBehaviour")
	}
}

type ElevatorState struct {
	Behaviour ElevatorBehaviour
	Floor     int
	Direction Direction
}

func (elevatorState ElevatorState) String() string {
	var builder strings.Builder

	builder.WriteString("ElevatorState{")
	builder.WriteString("\n\tBehaviour ")
	builder.WriteString(elevatorState.Behaviour.String())
	builder.WriteString("\n\tFloor ")
	fmt.Fprintf(&builder, "%d", elevatorState.Floor)
	builder.WriteString("\n\tDirection ")
	builder.WriteString(elevatorState.Direction.String())

	builder.WriteString("\n}")

	return builder.String()
}
