package elevator

import (
	"fmt"
	. "project/shared"
	"strings"
	"time"

	"github.com/angrycompany16/driver-go/elevio"
)

type ElevPositioning struct {
	direction            Direction
	behaviour            ElevatorBehaviour
	lastFloor            int
	floorBelow           int
	isAtFloor            bool
	door                 Door
	targetFloor          int // Will be -1 for no target
	enterFloor           chan int
	leaveFloor           chan int
	floorMovementTimeout *time.Timer
}

func InitPositioning() ElevPositioning {
	var door Door
	err := door.Close()
	if err != nil {
		panic("Failed to close the door at start up. This must be done to make sure the door is closed before the elevator starts to move")
	}

	motorStartTime := time.Now()
	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() == -1 {
		if time.Since(motorStartTime) > TimeBetweenFloors {
			panic("Failed to determine the elevator position within the expected time. The elevator is probably not working correctly")
		}
	}
	elevio.SetMotorDirection(elevio.MD_Stop)

	var pos ElevPositioning
	// pos.direction = UP
	// pos.behaviour = IDLE
	pos.lastFloor = elevio.GetFloor()
	pos.floorBelow = elevio.GetFloor()
	pos.isAtFloor = true
	pos.door = door
	pos.enterFloor = make(chan int)
	pos.leaveFloor = make(chan int)
	pos.floorMovementTimeout = time.NewTimer(TimeBetweenFloors)

	go pollPositionUpdates(pos.enterFloor, pos.leaveFloor)

	return pos
}

func pollPositionUpdates(enterFloor, leaveFloor chan<- int) {
	ticker := time.NewTicker(PositionPollInterval)
	defer ticker.Stop()

	lastFloor := elevio.GetFloor()

	for range ticker.C {
		floor := elevio.GetFloor()

		if lastFloor != floor {
			if floor == -1 {
				leaveFloor <- lastFloor
			} else {
				enterFloor <- floor
			}
			lastFloor = floor
		}
	}
}

func (pos *ElevPositioning) drive(direction Direction) {
	// if pos.behaviour == IDLE && pos.isAtFloor {
	// 	pos.lastSuccessState = time.Now()
	// }
	if (pos.door.IsOpen()) {
		panic("Cannot start to move the elevator while the door is open.")
	}

	// pos.behaviour = MOVING
	pos.direction = direction
	if direction == UP {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func (pos *ElevPositioning) stop() {
	// pos.behaviour = IDLE
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func (pos *ElevPositioning) String() string {
	var builder strings.Builder

	builder.WriteString("ElevatorController{\n")

	builder.WriteString(fmt.Sprintf("\tDirection: %s,\n", pos.direction.String()))
	builder.WriteString(fmt.Sprintf("\tBehaviour: %s,\n", pos.behaviour.String()))
	builder.WriteString(fmt.Sprintf("\tLastFloor: %d,\n", pos.lastFloor))
	builder.WriteString(fmt.Sprintf("\tFloorBelow: %d,\n", pos.floorBelow))
	builder.WriteString(fmt.Sprintf("\tIsAtFloor: %t,\n", pos.isAtFloor))
	builder.WriteString(fmt.Sprintf("\tTargetFloor: %d\n", pos.targetFloor))

	builder.WriteString("}")

	return builder.String()
}

func (pos *ElevPositioning) GetElevatorState() ElevatorState {
	var elevatorState ElevatorState

	elevatorState.Behaviour = pos.behaviour
	elevatorState.Direction = pos.direction
	elevatorState.Floor = pos.lastFloor

	return elevatorState
}

func (pos *ElevPositioning) preparePassengerTransfer() {
	if !pos.isAtFloor {
		panic("Cannot prepare Passenger transfer, with not at a floor.")
	}
	if pos.behaviour != IDLE {
		panic("Cannot prepare Passenger transfer, if the state is not IDLE.")
	}

	pos.behaviour = PASSENGER_TRANSFER
	fmt.Println("Opening door")
	pos.door.Open()

	if (pos.lastFloor == pos.targetFloor) {
		pos.targetFloor = -1
		// TODO: Tell the order system that we have stopped.
	}
}

// Warning: This function does not check if the elevator is in a faulty state
func (pos *ElevPositioning) driveToTarget() {
	if pos.isAtFloor && pos.targetFloor == pos.lastFloor {
		if pos.behaviour == MOVING {
			pos.stop()
			pos.behaviour = IDLE
		}
		pos.preparePassengerTransfer()
		return
	}

	if pos.targetFloor > pos.floorBelow {
		pos.drive(UP)
		pos.behaviour = MOVING
	}

	if pos.targetFloor <= pos.floorBelow {
		pos.drive(DOWN)
		pos.behaviour = MOVING
	}
}

func (pos *ElevPositioning) handleEnterFloor(floor int) {
	pos.lastFloor = floor
	pos.floorBelow = floor
	pos.isAtFloor = true

	// TODO: Maybe this should be moved to a own function, since this is a side effect.
	elevio.SetFloorIndicator(floor)

	switch pos.behaviour {
	case IDLE:
		fallthrough
	case PASSENGER_TRANSFER:
		fallthrough
	case DOOR_OBSTRUCTED:
		// TODO: This needs to be implemented
		panic("Enter faulty motor")
	case FAULTY_MOTOR:
		pos.behaviour = MOVING
		fallthrough
	case MOVING:
		if pos.targetFloor == floor {
			pos.stop()
			pos.behaviour = IDLE
			pos.preparePassengerTransfer()
		}
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (pos *ElevPositioning) handleLeaveFloor(floor int) {
	pos.isAtFloor = false
	if pos.direction == DOWN {
		pos.floorBelow = pos.lastFloor - 1
	}

	switch pos.behaviour {
	case IDLE:
		fallthrough
	case PASSENGER_TRANSFER:
		fallthrough
	case DOOR_OBSTRUCTED:
		// TODO: This needs to be implemented
		panic("Enter faulty motor")
	case FAULTY_MOTOR:
		pos.behaviour = MOVING
		fallthrough
	case MOVING:
		// This is normal
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (pos *ElevPositioning) handleTargetFloor(targetFloor int) {
	pos.targetFloor = targetFloor

	switch pos.behaviour {
	case IDLE:
		fallthrough
	case MOVING:
		pos.driveToTarget()
	case FAULTY_MOTOR:
		fallthrough
	case DOOR_OBSTRUCTED:
		fallthrough
	case PASSENGER_TRANSFER:
		// Do nothing
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (pos *ElevPositioning) handleCloseDoorTrigger() {
	switch pos.behaviour {
	case IDLE:
		fallthrough
	case FAULTY_MOTOR:
		fallthrough
	case MOVING:
		panic("The elevator got a CLOSE DOOR TRIGGER, but in the wrong state. The current state is " + pos.behaviour.String())
	case DOOR_OBSTRUCTED:
		fallthrough
	case PASSENGER_TRANSFER:
		fmt.Println("Trying to close door")
		err := pos.door.Close()
		if err != nil {
			fmt.Println("Failed to close door")
			pos.behaviour = DOOR_OBSTRUCTED
			pos.door.Open()
		} else {
			pos.behaviour = IDLE

			if pos.targetFloor == -1 {
				fmt.Println("Door closed, now idle")
			} else {
				fmt.Println("Door closed, driving to next target")
				pos.driveToTarget()
			}
		}
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (pos *ElevPositioning) handleDriving(targetFloor <-chan int) {
	for {
		select {
		case floor := <-pos.enterFloor:
			pos.handleEnterFloor(floor)
		case floor := <-pos.leaveFloor:
			pos.handleLeaveFloor(floor)
		case targetFloor := <-targetFloor:
			pos.handleTargetFloor(targetFloor)
		case <-pos.door.CloseTrigger():
			fmt.Println("Close door trigger")
			fmt.Println(pos.String())
			pos.handleCloseDoorTrigger()
		}
	}
}
