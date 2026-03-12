package elevator

import (
	"fmt"
	."project/orderHandler"
	. "project/shared"
	"github.com/angrycompany16/driver-go/elevio"
)

func ElevatorProcess(orderHandler *OrderHandler, elevatorStateCh chan<- ElevatorState, orderCompletedCh chan<- OrderCompleted, targetFloorCh <-chan int) {
	elevio.Init(elevatorServer, NumberOfFloors)
	elevio.SetStopLamp(false)
	positioning := InitPositioning(elevatorStateCh, orderCompletedCh)

	fmt.Println("Elevator state is determined.")

	go handleButtonPresses(orderHandler)
	go handleLights(orderHandler)

	positioning.handleDriving(targetFloorCh)
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
