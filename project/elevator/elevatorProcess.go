package elevator

import (
	"github.com/angrycompany16/driver-go/elevio"
	. "project/shared"
)

func ElevatorProcess(channels ElevatorInterface, elevatorServerHost string) {
	elevio.Init(elevatorServerHost, NumberOfFloors)
	elevio.SetStopLamp(false)

	go handleButtonPresses(channels.NewOrderCh)
	go handleLights(channels.ConfirmedOrdersCh)

	startElevatorController(channels.ElevatorStateCh, channels.OrderCompletedCh, channels.TargetFloorCh)
}
