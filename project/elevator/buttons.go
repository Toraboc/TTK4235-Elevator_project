package elevator

import (
	"fmt"
	. "project/shared"
	. "project/orderHandler"
	"time"

	"github.com/angrycompany16/driver-go/elevio"
)

func checkForButtonPress(newOrderCh chan<- NewOrderEvent, floor int, buttonType elevio.ButtonType, orderType OrderType, floorButtonState *[NumberOfFloors][3]bool) {
	newValue := elevio.GetButton(buttonType, floor)
	oldValue := floorButtonState[floor][orderType]

	if newValue && !oldValue {
		fmt.Printf("Button pressed floor = %d, type = %v\n", floor, orderType)
		newOrderCh <- NewOrderEvent{Floor: floor, OrderType: orderType}
	}

	floorButtonState[floor][orderType] = newValue
}

func handleButtonPresses(newOrderCh chan<- NewOrderEvent) {
	var floorButtonState [NumberOfFloors][3]bool
	for {
		time.Sleep(40 * time.Millisecond)
		for i := range NumberOfFloors {
			checkForButtonPress(newOrderCh, i, elevio.BT_HallUp, HALLUP, &floorButtonState)
			checkForButtonPress(newOrderCh, i, elevio.BT_HallDown, HALLDOWN, &floorButtonState)
			checkForButtonPress(newOrderCh, i, elevio.BT_Cab, CAB, &floorButtonState)
		}
	}
}
