package elevator

import (
	"Driver-go/elevio"
	. "project/shared"
	"fmt"
)

func ElevatorProcess() {
	elevio.Init("localhost:15657", 4)
	positioning := InitPositioning()

	fmt.Println("Elevator state is determined.")

	go handleButtonPresses()

	positioning.handleDriving()
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
