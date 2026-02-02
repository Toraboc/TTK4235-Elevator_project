package elevator

import (
	"Driver-go/elevio"
	. "project/shared"
	"time"
	"fmt"
)

func getButtonsPresses() ([NumberOfFloors][3]bool, bool) {
	var buttonPresses [NumberOfFloors][3]bool
	anyPressed := false
	for i := range NumberOfFloors {
		buttonPresses[i][0] = elevio.GetButton(elevio.BT_HallUp, i)
		buttonPresses[i][1] = elevio.GetButton(elevio.BT_HallDown, i)
		buttonPresses[i][2] = elevio.GetButton(elevio.BT_Cab, i)
		anyPressed = anyPressed || buttonPresses[i][0] || buttonPresses[i][1] || buttonPresses[i][2]
	}
	return buttonPresses, anyPressed
}

func handleButtonPresses() {
	for {
		time.Sleep(20 * time.Millisecond)
		_, anyPressed := getButtonsPresses()
		if (anyPressed) {
			fmt.Println("Some buttons are pressed")
			// TODO: Notify the orderHandler
		}
	}
}
