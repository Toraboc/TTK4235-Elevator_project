package main

import (
	"Driver-go/elevio"
)

type Direction int

const (
	DirUp Direction  = iota
	DirDown
	DirStop
)

type Position struct {
	direction Direction
	lastDirection Direction
	lastFloor int
	floorBelow int
	isAtAFloor bool
	targetFloor int
}


func getPosition(pos *Position) int {
	return pos.lastFloor
}

func getDirection(pos *Position) Direction {
	return pos.direction
}

func gotoFloor(pos *Position, door *Door, floor int) {
	pos.targetFloor = floor
	door.willOpenDoor = true
}

func initPosition(pos *Position) {
	pos.direction = DirStop
	pos.lastDirection = DirDown
	pos.lastFloor = 0
	pos.targetFloor = -1

	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() == -1 {}
	pos.targetFloor = 0
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func positionModuleLoop(pos *Position, door *Door) {
	if elevio.GetFloor() != -1 {
		pos.lastFloor = elevio.GetFloor()
		pos.floorBelow = pos.lastFloor
		pos.isAtAFloor = true
		elevio.SetFloorIndicator(pos.lastFloor)
	} else {
		if pos.isAtAFloor {
			if pos.direction == DirDown {
				pos.floorBelow = pos.lastFloor - 1
			}
		}
		pos.isAtAFloor = false
	}
	if pos.targetFloor == -1 {
		pos.direction = DirStop
		elevio.SetMotorDirection(elevio.MD_Stop)
	} else if pos.isAtAFloor && pos.lastFloor == pos.targetFloor {
		// Elevator has arrived at target floor
		pos.direction = DirStop
		elevio.SetMotorDirection(elevio.MD_Stop)
		pos.targetFloor = -1
		if getDoorState(door) == DoorClosed && door.willOpenDoor {
			door.willOpenDoor = false
			openDoor(door)
		}
	} else if pos.floorBelow < pos.targetFloor {
		// Elevator needs to go up
		if getDoorState(door) == DoorClosed {
			pos.direction = DirUp
			pos.lastDirection = DirUp
			elevio.SetMotorDirection(elevio.MD_Up)
		}
	} else if pos.floorBelow >= pos.targetFloor {
		// Elevator needs to go down
		if getDoorState(door) == DoorClosed {
			pos.direction = DirDown
			pos.lastDirection = DirDown
			elevio.SetMotorDirection(elevio.MD_Down)
		}
	}
}