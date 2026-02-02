package elevator

import (
	"Driver-go/elevio"
	. "project/shared"
	"time"
	"fmt"
)

type ElevPositioning struct {
	direction             Direction
	behaviour             ElevatorBehaviour
	lastFloorSensorChange time.Time
	lastFloor             int
	floorBelow            int
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
	pos.isAtFloor = true
	pos.door = door

	return pos
}

func (pos *ElevPositioning) updatePosition() {
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

func (pos *ElevPositioning) drive(direction Direction) {
	pos.behaviour = MOVING
	pos.direction = direction
	if direction == UP {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func (pos *ElevPositioning) stop() {
	pos.behaviour = IDLE
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func (pos *ElevPositioning) handleElevatorMotor(targetFloor int) {

	// TODO: this needs some cleanup
	if targetFloor == -1 {
		if (pos.behaviour == MOVING) {
			pos.stop()
		}
	} else if pos.isAtFloor && pos.lastFloor == targetFloor {
		// We are at the target
		// Check if we just arrived
		if (pos.behaviour == MOVING) {
			pos.stop()

			pos.door.Open()
			pos.behaviour = PASSENGER_TRANSFER
			// TODO: Tell the orderHandler that we have stopped
		}

	} else if (pos.behaviour == IDLE || pos.behaviour == MOVING) {
		if pos.floorBelow < targetFloor {
			pos.drive(UP)
		} else if (pos.floorBelow >= targetFloor) {
			pos.drive(DOWN)
		}
	}
}

func (pos *ElevPositioning) printState() {
	fmt.Println("lastFloor:", pos.lastFloor, "floorBelow", pos.floorBelow, "direction", pos.direction, "behaviour", pos.behaviour)
}

func (pos *ElevPositioning) handleDriving() {
	for {
		time.Sleep(50 * time.Millisecond)
		pos.updatePosition()

		// TODO: Change the target floor to a floor from the orderHandler
		pos.handleElevatorMotor(0)


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
