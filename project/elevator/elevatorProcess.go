package elevator

import (
	"github.com/angrycompany16/driver-go/elevio"
	. "project/orderHandler"
	. "project/shared"
)

func ElevatorProcess(orderHandler *OrderHandler, elevatorStateCh chan<- ElevatorState, orderCompletedCh chan<- OrderCompleted, targetFloorCh <-chan int, orderNewCh chan<- OrderNew) {
	elevio.Init(elevatorServer, NumberOfFloors)
	elevio.SetStopLamp(false)

	go handleButtonPresses(orderNewCh)
	go handleLights(orderHandler)

	startElevatorController(elevatorStateCh, orderCompletedCh, targetFloorCh)
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
