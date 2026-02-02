package elevator

import (
	"Driver-go/elevio"
	. "project/shared"
	"time"
)

type ElevPositioning struct {
	direction             Direction
	behaviour             ElevatorBehaviour
	lastFloorSensorChange time.Time
	lastFloor             int
	floorBelow            int
	targetFloor           int
	isAtFloor             bool
	door                  Door
}

func InitPositioning() ElevPositioning {
	var door Door
	err := door.Close()
	if err != nil {
		panic("Failed to close the door at start up. This must be done to make sure the door is closed before the elevator starts to move")
	}

	// TODO: The elevator could be broken at this point, some checks should be implemented
	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() == -1 {
	}
	elevio.SetMotorDirection(elevio.MD_Stop)

	var pos ElevPositioning
	pos.direction = UP
	pos.behaviour = IDLE
	pos.lastFloorSensorChange = time.Now()
	pos.lastFloor = elevio.GetFloor()
	pos.floorBelow = elevio.GetFloor()
	pos.targetFloor = -1
	pos.isAtFloor = true
	pos.door = door

	return pos
}

func (pos ElevPositioning) updatePosition() {
	// TODO: This logic does not handle if someone moves the elevator by force
	// We could implemented a check that the floorBelow only is updated if the behvaiour is MOVING, else panic or something
	// or go into an obstructed state
	floor := elevio.GetFloor()
	if floor != -1 {
		if !pos.isAtFloor {
			pos.lastFloorSensorChange = time.Now()
		}
		pos.lastFloor = floor
		pos.floorBelow = floor
		pos.isAtFloor = true
		elevio.SetFloorIndicator(pos.lastFloor)
	} else if pos.isAtFloor {
		pos.isAtFloor = false
		pos.lastFloorSensorChange = time.Now()
		if pos.direction == DOWN {
			pos.floorBelow--
		}
	}
}

func (pos ElevPositioning) drive(direction Direction) {
	pos.behaviour = MOVING
	pos.direction = direction
	if direction == UP {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func (pos ElevPositioning) stop() {
	pos.behaviour = IDLE
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func (pos ElevPositioning) handleElevatorMotor() {

	// TODO: this needs some cleanup
	if pos.targetFloor == -1 {
		if (pos.behaviour == MOVING) {
			pos.stop()
		}
	} else if pos.isAtFloor && pos.lastFloor == pos.targetFloor {
		// We are at the target
		// Check if we just arrived
		if (pos.behaviour == MOVING) {
			pos.stop()

			// TODO: Tell the orderHandler that we have stopped

			pos.door.Open()
			pos.behaviour = PASSENGER_TRANSFER
		}

	} else if (pos.behaviour == IDLE) {
		if pos.floorBelow <= pos.targetFloor {
			pos.drive(DOWN)
		} else if (pos.floorBelow > pos.targetFloor) {
			pos.drive(UP)
		}
	}
}

func (pos ElevPositioning) handleDriving() {
	for {
		pos.updatePosition()
		pos.handleElevatorMotor()

		// Close the door after some time
		if (pos.behaviour == PASSENGER_TRANSFER || (pos.behaviour == OBSTRCTED && pos.door.IsOpen())) {
			if (time.Since(pos.door.changeTime) > doorOpenTime) {
				err := pos.door.Close()
				if err != nil {
					pos.behaviour = OBSTRCTED
				} else {
					pos.behaviour = IDLE
				}
			}
		}

		// TODO: Add some logic to check if the elevator uses too much time on movements
	}

}
