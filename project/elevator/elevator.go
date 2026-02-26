package elevator

import (
	"fmt"
	."project/orderHandler"
	. "project/shared"
	"github.com/angrycompany16/driver-go/elevio"
)

func ElevatorProcess(orderHandler *OrderHandler) {
	elevio.Init("localhost:15657", 4)
	positioning := InitPositioning()

	fmt.Println("Elevator state is determined.")

	elevio.SetStopLamp(false)

	go handleButtonPresses()
	go handleLights(orderHandler)

	positioning.handleDriving(orderHandler)
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
