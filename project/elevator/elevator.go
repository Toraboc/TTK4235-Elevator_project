package elevator

import (
	"fmt"
	"Driver-go/elevio"
	. "project/shared"
)

func ElevatorProcess() {
	elevio.Init("localhost:15657", 4)
	positioning := InitPositioning()

	fmt.Println("Elevator state is determined.")

	go handleButtonPresses()
	go handleLights()

	positioning.handleDriving()
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
