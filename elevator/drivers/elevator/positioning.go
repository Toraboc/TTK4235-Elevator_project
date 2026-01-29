package elevator

import (
	"Driver-go/elevio"
)

type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirStop
)

type Position struct {
	direction     Direction
	LastDirection Direction
	LastFloor     int
	floorBelow    int
	isAtAFloor    bool
	TargetFloor   int
}

/*
getPosition(pointer to Position) int
Returns the last known floor of the elevator.
*/
func getPosition(pos *Position) int {
	return pos.LastFloor
}

/*
getDirection(pointer to Position) Direction
Returns the current direction of the elevator.
*/
func GetDirection(pos *Position) Direction {
	return pos.direction
}

func GotoFloor(pos *Position, door *Door, floor int) {
	pos.TargetFloor = floor
	door.WillOpen = true
}

func InitPosition(pos *Position) {
	pos.direction = DirStop
	pos.LastDirection = DirDown
	pos.LastFloor = 0
	pos.TargetFloor = -1

	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() == -1 {
	}
	pos.TargetFloor = 0
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func PositionModuleLoop(pos *Position, door *Door) {
	temp_floor := elevio.GetFloor()
	if temp_floor != -1 { // Elevator is at a floor
		pos.LastFloor = temp_floor
		pos.floorBelow = pos.LastFloor
		pos.isAtAFloor = true
		elevio.SetFloorIndicator(pos.LastFloor)
	} else {
		if pos.isAtAFloor {
			if pos.direction == DirDown {
				pos.floorBelow = pos.LastFloor - 1
			}
		}
		pos.isAtAFloor = false
	}

	if pos.TargetFloor == -1 {
		pos.direction = DirStop
		elevio.SetMotorDirection(elevio.MD_Stop)
	} else if pos.isAtAFloor && pos.LastFloor == pos.TargetFloor {
		// Elevator has arrived at target floor
		pos.direction = DirStop
		elevio.SetMotorDirection(elevio.MD_Stop)
		pos.TargetFloor = -1
		if !(door.IsOpen) && door.WillOpen {
			door.WillOpen = false
			openDoor(door)
		}
	} else if pos.floorBelow < pos.TargetFloor {
		// Elevator needs to go up
		if !(door.IsOpen) {
			pos.direction = DirUp
			pos.LastDirection = DirUp
			elevio.SetMotorDirection(elevio.MD_Up)
		}
	} else if pos.floorBelow >= pos.TargetFloor {
		// Elevator needs to go down
		if !(door.IsOpen) {
			pos.direction = DirDown
			pos.LastDirection = DirDown
			elevio.SetMotorDirection(elevio.MD_Down)
		}
	}
}
