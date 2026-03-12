package elevator

import (
	"fmt"
	. "project/shared"
	"time"

	"github.com/angrycompany16/driver-go/elevio"
)

func checkForButtonPress(orderNewCh chan<- OrderNew, floor int, buttonType elevio.ButtonType, orderType OrderType, floorButtonState *[NumberOfFloors][3]bool) {
	newValue := elevio.GetButton(buttonType, floor)
	oldValue := floorButtonState[floor][orderType]

	if newValue && !oldValue {
		fmt.Printf("Button pressed floor = %d, type = %v\n", floor, orderType)
		orderNewCh <- OrderNew{Floor: floor, Type: orderType}
	}

	floorButtonState[floor][orderType] = newValue
}

func handleButtonPresses(orderNewCh chan<- OrderNew) {
	var floorButtonState [NumberOfFloors][3]bool
	for {
		time.Sleep(40 * time.Millisecond)
		for i := range NumberOfFloors {
			checkForButtonPress(orderNewCh, i, elevio.BT_HallUp, HALLUP, &floorButtonState)
			checkForButtonPress(orderNewCh, i, elevio.BT_HallDown, HALLDOWN, &floorButtonState)
			checkForButtonPress(orderNewCh, i, elevio.BT_Cab, CAB, &floorButtonState)
		}
	}
}

