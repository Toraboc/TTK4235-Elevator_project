package elevator

import (
	"github.com/angrycompany16/driver-go/elevio"
	. "project/orderHandler"
	. "project/shared"
)

func ElevatorProcess(channels OrderChannels) {
	elevio.Init(elevatorServer, NumberOfFloors)
	elevio.SetStopLamp(false)

	go handleButtonPresses(channels)
	go handleLights(channels)

	startElevatorController(channels.ElevatorStateCh, channels.OrderCompletedCh, channels.TargetFloorCh)
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
