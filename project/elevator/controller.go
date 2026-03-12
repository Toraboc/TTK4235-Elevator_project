package elevator

import (
	"fmt"
	. "project/shared"
	"strings"
	"time"

	"github.com/angrycompany16/driver-go/elevio"
)

type ElevatorDetailedState struct {
	direction   Direction
	behaviour   ElevatorBehaviour
	lastFloor   int
	floorBelow  int
	isAtFloor   bool
	targetFloor int // Will be -1 for no target
}

type ElevatorController struct {
	state                ElevatorDetailedState
	door                 Door
	elevatorStateCh      chan<- ElevatorState
	orderCompletedCh     chan<- OrderCompleted
	floorMovementTimeout *time.Timer
	lastElevatorState    ElevatorState
}

func startElevatorController(elevatorStateCh chan<- ElevatorState, orderCompletedCh chan<- OrderCompleted, targetFloorCh <-chan int) {
	var door Door
	err := door.Close()
	if err != nil {
		panic("Failed to close the door at start up. This must be done to make sure the door is closed before the elevator starts to move")
	}

	motorStartTime := time.Now()
	elevio.SetMotorDirection(elevio.MD_Down)
	for elevio.GetFloor() == -1 {
		if time.Since(motorStartTime) > timeBetweenFloors {
			panic("Failed to determine the elevator position within the expected time. The elevator is probably not working correctly")
		}
	}
	elevio.SetMotorDirection(elevio.MD_Stop)

	var detailedState ElevatorDetailedState
	detailedState.behaviour = IDLE
	detailedState.lastFloor = elevio.GetFloor()
	detailedState.floorBelow = elevio.GetFloor()
	detailedState.isAtFloor = true
	detailedState.targetFloor = -1

	var controller ElevatorController
	controller.state = detailedState
	controller.door = door
	controller.elevatorStateCh = elevatorStateCh
	controller.orderCompletedCh = orderCompletedCh
	controller.floorMovementTimeout = time.NewTimer(timeBetweenFloors)
	controller.floorMovementTimeout.Stop()

	enterFloorCh := make(chan int)
	leaveFloorCh := make(chan int)

	go pollPositionUpdates(enterFloorCh, leaveFloorCh)

	controller.lastElevatorState = controller.GetElevatorState()
	controller.elevatorStateCh <- controller.lastElevatorState

	fmt.Println("Elevator state is determined.")

	controller.handleDriving(targetFloorCh, enterFloorCh, leaveFloorCh)
}

func (controller *ElevatorController) sendElevatorStateUpdate() {
	newElevatorState := controller.GetElevatorState()

	if newElevatorState.Behaviour != controller.lastElevatorState.Behaviour ||
		newElevatorState.Direction != controller.lastElevatorState.Direction ||
		newElevatorState.Floor != controller.lastElevatorState.Floor {
		controller.lastElevatorState = newElevatorState
		controller.elevatorStateCh <- newElevatorState
	}
}

func pollPositionUpdates(enterFloor, leaveFloor chan<- int) {
	ticker := time.NewTicker(positionPollInterval)
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

func (controller *ElevatorController) drive(direction Direction) {
	// if pos.behaviour == IDLE && pos.isAtFloor {
	// 	pos.lastSuccessState = time.Now()
	// }
	if controller.door.IsOpen() {
		panic("Cannot start to move the elevator while the door is open.")
	}

	// pos.behaviour = MOVING
	controller.state.direction = direction
	if direction == UP {
		elevio.SetMotorDirection(elevio.MD_Up)
	} else {
		elevio.SetMotorDirection(elevio.MD_Down)
	}
}

func (controller *ElevatorController) stop() {
	// pos.behaviour = IDLE
	elevio.SetMotorDirection(elevio.MD_Stop)
}

func (elevator *ElevatorDetailedState) String() string {
	var builder strings.Builder

	builder.WriteString("ElevatorDetailedState {\n")

	fmt.Fprintf(&builder, "\tDirection: %v,\n", elevator.direction)
	fmt.Fprintf(&builder, "\tBehaviour: %v,\n", elevator.behaviour)
	fmt.Fprintf(&builder, "\tLastFloor: %d,\n", elevator.lastFloor)
	fmt.Fprintf(&builder, "\tFloorBelow: %d,\n", elevator.floorBelow)
	fmt.Fprintf(&builder, "\tIsAtFloor: %t,\n", elevator.isAtFloor)
	fmt.Fprintf(&builder, "\tTargetFloor: %d,\n", elevator.targetFloor)

	builder.WriteString("}")

	return builder.String()
}

func (controller *ElevatorController) String() string {
	var builder strings.Builder

	builder.WriteString("ElevatorController {\n")

	fmt.Fprintf(&builder, "\tstate: %s,\n", strings.ReplaceAll(controller.state.String(), "\n", "\n\t"))
	fmt.Fprintf(&builder, "\tdoorOpen: %t\n", controller.door.IsOpen())

	builder.WriteString("}")

	return builder.String()

}

func (controller *ElevatorController) GetElevatorState() ElevatorState {
	var elevatorState ElevatorState

	elevatorState.Behaviour = controller.state.behaviour
	elevatorState.Direction = controller.state.direction
	elevatorState.Floor = controller.state.lastFloor

	return elevatorState
}

func (controller *ElevatorController) preparePassengerTransfer() {
	if !controller.state.isAtFloor {
		panic("Cannot prepare Passenger transfer, with not at a floor.")
	}
	if controller.state.behaviour != IDLE && controller.state.behaviour != PASSENGER_TRANSFER {
		panic("Cannot prepare Passenger transfer, if the state is not IDLE or PASSENGER_TRANSFER.")
	}

	controller.state.behaviour = PASSENGER_TRANSFER
	controller.door.Open()
	controller.orderCompletedCh <- OrderCompleted{Floor: controller.state.lastFloor, Direction: controller.state.direction}
	if (controller.state.lastFloor == controller.state.targetFloor) {
		controller.state.targetFloor = -1
	}
}

// Warning: This function does not check if the elevator is in a faulty state
func (controller *ElevatorController) driveToTarget() {
	if controller.state.isAtFloor && controller.state.targetFloor == controller.state.lastFloor {
		if controller.state.behaviour == MOVING {
			controller.stop()
			controller.state.behaviour = IDLE
		}
		controller.preparePassengerTransfer()
		return
	}

	if controller.state.targetFloor == -1 {
		return
	}

	if controller.state.targetFloor > controller.state.floorBelow {
		controller.drive(UP)
		controller.state.behaviour = MOVING
	}

	if controller.state.targetFloor <= controller.state.floorBelow {
		controller.drive(DOWN)
		controller.state.behaviour = MOVING
	}
}

func (controller *ElevatorController) handleEnterFloor(floor int) {
	switch controller.state.behaviour {
	case IDLE:
		fallthrough
	case PASSENGER_TRANSFER:
		fallthrough
	case DOOR_OBSTRUCTED:
		panic("Cannot enter a floor in the state " + controller.state.behaviour.String() + ".")
	case FAULTY_MOTOR:
		controller.state.behaviour = MOVING
		fallthrough
	case MOVING:
		expectedFloorDelta := -1
		if (controller.state.direction == UP) {
			expectedFloorDelta = 1
		}
		expectedFloor := controller.state.lastFloor + expectedFloorDelta
		if expectedFloor != floor {
			panic(fmt.Sprintf("The elevator reached the floor %d, but expected to reach floor %d. Something is terrably wrong.", floor, expectedFloor))
		}

		controller.state.lastFloor = floor
		controller.state.floorBelow = floor
		controller.state.isAtFloor = true
		controller.floorMovementTimeout.Stop()

		// TODO: Maybe this should be moved to a own function, since this is a side effect.
		elevio.SetFloorIndicator(floor)

		if controller.state.targetFloor == -1 {
			controller.stop()
			controller.state.behaviour = IDLE
		}

		if controller.state.targetFloor == floor {
			controller.stop()
			controller.state.behaviour = IDLE
			controller.preparePassengerTransfer()
		}
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}

}

func (controller *ElevatorController) handleLeaveFloor(_ int) {
	controller.state.isAtFloor = false
	if controller.state.direction == DOWN {
		controller.state.floorBelow = controller.state.lastFloor - 1
	}

	controller.floorMovementTimeout.Reset(timeBetweenFloors)

	switch controller.state.behaviour {
	case IDLE:
		controller.state.behaviour = FAULTY_MOTOR
		if controller.state.lastFloor < 2 {
			controller.drive(UP)
		} else {
			controller.drive(DOWN)
		}
	case PASSENGER_TRANSFER:
		fallthrough
	case DOOR_OBSTRUCTED:
		panic("The elevator left the floor, and the door is open, it's a bit unclear what should happen.")
	case FAULTY_MOTOR:
		controller.state.behaviour = MOVING
		fallthrough
	case MOVING:
		// This is normal
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (controller *ElevatorController) handleTargetFloor(targetFloor int) {
	controller.state.targetFloor = targetFloor

	switch controller.state.behaviour {
	case IDLE:
		fallthrough
	case MOVING:
		controller.driveToTarget()
	case FAULTY_MOTOR:
		// Do nothing
	case DOOR_OBSTRUCTED:
		fallthrough
	case PASSENGER_TRANSFER:
		controller.preparePassengerTransfer()
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (controller *ElevatorController) handleCloseDoorTrigger() {
	switch controller.state.behaviour {
	case IDLE:
		fallthrough
	case FAULTY_MOTOR:
		fallthrough
	case MOVING:
		panic("The elevator got a CLOSE DOOR TRIGGER, but in the wrong state. The current state is " + controller.state.behaviour.String())
	case DOOR_OBSTRUCTED:
		fallthrough
	case PASSENGER_TRANSFER:
		err := controller.door.Close()
		if err != nil {
			controller.state.behaviour = DOOR_OBSTRUCTED
			controller.door.Open()
		} else {
			controller.state.behaviour = IDLE

			if controller.state.targetFloor != -1 {
				controller.driveToTarget()
			}
		}
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (controller *ElevatorController) handleFloorMovementTimeout() {
	switch controller.state.behaviour {
	case IDLE:
		fallthrough
	case DOOR_OBSTRUCTED:
		fallthrough
	case FAULTY_MOTOR:
		fallthrough
	case PASSENGER_TRANSFER:
		panic("The floor movement timeout triggered, when we are not in the MOVING state, this should never happen.")
	case MOVING:
		controller.state.behaviour = FAULTY_MOTOR
	case DISCONNECTED:
		panic("Our elevator can never become DISCONNECTED")
	}
}

func (controller *ElevatorController) handleDriving(targetFloorCh, enterFloorCh, leaveFloorCh <-chan int) {
	for {
		fmt.Println("Listening for events...")
		select {
		case floor := <-enterFloorCh:
			fmt.Printf("ENTER FLOOR %d\n", floor)
			controller.handleEnterFloor(floor)
		case floor := <-leaveFloorCh:
			fmt.Printf("LEAVE FLOOR %d\n", floor)
			controller.handleLeaveFloor(floor)
		case targetFloor := <-targetFloorCh:
			fmt.Printf("TARGET FLOOR = %d\n", targetFloor)
			controller.handleTargetFloor(targetFloor)
		case <-controller.door.CloseTrigger():
			fmt.Println("CLOSE DOOR")
			controller.handleCloseDoorTrigger()
		case <-controller.floorMovementTimeout.C:
			fmt.Println("FLOOR MOVEMENT TIMEOUT")
			controller.handleFloorMovementTimeout()
		}
		fmt.Printf("AFTER EVENT STATE:\n%v\n", controller)
		controller.sendElevatorStateUpdate()
	}
}
