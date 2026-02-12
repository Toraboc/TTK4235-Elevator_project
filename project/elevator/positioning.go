package elevator

import (
	"Driver-go/elevio"
	"fmt"
	. "project/shared"
	"time"
)

type ElevPositioning struct {
	direction        Direction
	behaviour        ElevatorBehaviour
	lastSuccessState time.Time
	lastFloor        int
	floorBelow       int
	isAtFloor        bool
	door             Door
}

const TimeBetweenFloors = 4 * time.Second

func InitPositioning() ElevPositioning {
	var door Door
	err := door.Close()
	if err != nil {
		panic("Failed to close the door at start up. This must be done to make sure the door is closed before the elevator starts to move")
	}

	motorStartTime := time.Now()
	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() == -1 {
		if (time.Since(motorStartTime) > TimeBetweenFloors) {
			panic("Failed to determine the elevator position within the expected time. The elevator is probably not working correctly")
		}
	}
	elevio.SetMotorDirection(elevio.MD_Stop)

	var pos ElevPositioning
	pos.direction = UP
	pos.behaviour = IDLE
	pos.lastSuccessState = time.Now()
	pos.lastFloor = elevio.GetFloor()
	pos.floorBelow = elevio.GetFloor()
	pos.isAtFloor = true
	pos.door = door

	return pos
}

func (pos *ElevPositioning) updatePosition() {
	// TODO: This logic does not handle if someone moves the elevator by force
	// We could implemented a check that the floorBelow only is updated if the behaviour is MOVING, else panic or something
	// or go into an obstructed state
	floor := elevio.GetFloor()
	if floor != -1 {
		if !pos.isAtFloor {
			// The elevator arrived at a floor
			pos.lastSuccessState = time.Now()
		}
		pos.lastFloor = floor
		pos.floorBelow = floor
		pos.isAtFloor = true
		elevio.SetFloorIndicator(pos.lastFloor)
	} else if pos.isAtFloor {
		// The elevator is leaving a floor
		pos.isAtFloor = false
		pos.lastSuccessState = time.Now()
		if pos.direction == DOWN {
			pos.floorBelow--
		}

		if (!pos.behaviour.CanMove()) {
			// This is really bad, the elevator shouold not be moving.
			// And we also don't know which floor the elevator is at.
			// Entry a faulty state
			if pos.door.IsOpen() {
				// TODO: Fix this panic
				// I'm not sure that the expected behaviour is in this case
				panic("The elevator was moving when it should be standing still. And the door is open, therefore we will not try to recover.")
			}

			pos.behaviour = FAULTY_MOTOR
			pos.lastSuccessState = time.Now().Add(-TimeBetweenFloors)

			// Find the most safe direction to try to drive in
			if (pos.lastFloor > 1) {
				pos.direction = DOWN
			} else {
				pos.direction = UP
			}
		}
	}

	// Detect if the elevator is using too much time on the movements
	if pos.behaviour == MOVING && time.Since(pos.lastSuccessState) > TimeBetweenFloors {
		pos.behaviour = FAULTY_MOTOR
	} else if (pos.behaviour == FAULTY_MOTOR && time.Since(pos.lastSuccessState) < TimeBetweenFloors) {
		pos.behaviour = MOVING
	}
}

func (pos *ElevPositioning) drive(direction Direction) {
	if (pos.behaviour == IDLE && pos.isAtFloor) {
		pos.lastSuccessState = time.Now()
	}
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

	fmt.Println("TargetFloor", targetFloor)

	if targetFloor == -1 {
		if pos.behaviour == MOVING && pos.isAtFloor {
			pos.stop()
		}
	} else if pos.isAtFloor && pos.lastFloor == targetFloor {
		// We are at the target
		// Check if we just arrived
		if pos.behaviour == MOVING {
			pos.stop()

			pos.door.Open()
			pos.behaviour = PASSENGER_TRANSFER
			// TODO: Tell the orderHandler that we have stopped
		}

	} else if pos.behaviour == IDLE || pos.behaviour == MOVING {
		if pos.floorBelow < targetFloor {
			pos.drive(UP)
		} else if pos.floorBelow >= targetFloor {
			pos.drive(DOWN)
		}
	}
}

func (pos *ElevPositioning) recoverFromFaultyMotor() {
	if (pos.behaviour != FAULTY_MOTOR) {
		panic("Cannot recover from a state that is not faulty motors.")
	}

	// Try to make the elevator continue in the last direction
	if pos.direction == UP {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func getBehaviourString(behaviour ElevatorBehaviour) string {
	switch (behaviour) {
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
		panic("Undefined behaviour")
	}
}

func (pos *ElevPositioning) printState() {
	fmt.Println("lastFloor:", pos.lastFloor, "floorBelow", pos.floorBelow, "direction", pos.direction, "behaviour", getBehaviourString(pos.behaviour))
}

func (pos *ElevPositioning) handleDriving() {
	for {
		time.Sleep(50 * time.Millisecond)
		pos.updatePosition()

		// TODO: Change the target floor to a floor from the orderHandler
		if (pos.behaviour == FAULTY_MOTOR) {
			pos.recoverFromFaultyMotor()
		} else if (pos.behaviour.CanBeAssignedOrders()) {
			pos.handleElevatorMotor(getAButtonFloor())
		}

		pos.printState()

		// Close the door after some time
		if pos.behaviour == PASSENGER_TRANSFER || pos.behaviour == DOOR_OBSTRUCTED {
			if time.Since(pos.door.changeTime) > doorOpenTime {
				err := pos.door.Close()
				if err != nil {
					pos.behaviour = DOOR_OBSTRUCTED
				} else {
					pos.behaviour = IDLE
				}
			}
		}
	}

}
