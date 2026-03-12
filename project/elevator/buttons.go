package elevator

import (
	"fmt"
	. "project/orderHandler"
	. "project/shared"
	"time"

	"github.com/angrycompany16/driver-go/elevio"
)

func checkForButtonPress(orderHandler *OrderHandler, floor int, buttonType elevio.ButtonType, orderType OrderType, floorButtonState *[NumberOfFloors][3]bool) {
	newValue := elevio.GetButton(buttonType, floor)
	oldValue := floorButtonState[floor][orderType]

	if newValue && !oldValue {
		fmt.Printf("Button pressed floor = %d, type = %v\n", floor, orderType)
		orderHandler.UpdateNewOrder(floor, orderType)
		floorButtonState[floor][orderType] = newValue
	}
}

func handleButtonPresses(orderHandler *OrderHandler) {
	var floorButtonState [NumberOfFloors][3]bool
	for {
		time.Sleep(40 * time.Millisecond)
		for i := range NumberOfFloors {
			checkForButtonPress(orderHandler, i, elevio.BT_HallUp, HALLUP, &floorButtonState)
			checkForButtonPress(orderHandler, i, elevio.BT_HallDown, HALLDOWN, &floorButtonState)
			checkForButtonPress(orderHandler, i, elevio.BT_Cab, CAB, &floorButtonState)
		}
	}
}

