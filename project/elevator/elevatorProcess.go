package elevator

import (
	"fmt"
	."project/orderHandler"
	. "project/shared"
	"github.com/angrycompany16/driver-go/elevio"
)

func ElevatorProcess(orderHandler *OrderHandler, targetFloorCh <-chan int) {
	elevio.Init(ElevatorServer, NumberOfFloors)
	elevio.SetStopLamp(false)
	positioning := InitPositioning()

	fmt.Println("Elevator state is determined.")

	go handleButtonPresses(orderHandler)
	go handleLights(orderHandler)

	positioning.handleDriving(targetFloorCh)
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
