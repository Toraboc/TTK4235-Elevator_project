package elevator

import (
	"fmt"
	. "project/shared"
	"github.com/angrycompany16/driver-go/elevio"
)

func ElevatorProcess() {
	elevio.Init("localhost:15657", 4)
	positioning := InitPositioning()

	fmt.Println("Elevator state is determined.")

	elevio.SetStopLamp(false)

	go handleButtonPresses()
	go handleLights()

	positioning.handleDriving()
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
