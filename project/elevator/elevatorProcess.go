package elevator

import (
	"github.com/angrycompany16/driver-go/elevio"
	. "project/orderHandler"
	. "project/shared"
)

func ElevatorProcess(elevatorServerHost string, channels OrderChannels) {
	elevio.Init(elevatorServerHost, NumberOfFloors)
	elevio.SetStopLamp(false)

	go handleButtonPresses(channels.NewOrderCh)
	go handleLights(channels.ConfirmedOrdersCh)

	startElevatorController(channels.ElevatorStateCh, channels.OrderCompletedCh, channels.TargetFloorCh)
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
