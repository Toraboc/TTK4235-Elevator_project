package elevator

import (
	. "project/orderHandler"
	. "project/shared"
	"time"

	"github.com/angrycompany16/driver-go/elevio"
)

func checkForButtonPress(orderHandler *OrderHandler, floor int, buttonType elevio.ButtonType, orderType OrderType) {
	if (elevio.GetButton(buttonType, floor)) {
		orderHandler.UpdateOrder(floor, orderType)
	}
}

func handleButtonPresses(orderHandler *OrderHandler) {
	for {
		time.Sleep(40 * time.Millisecond)
		for i := range NumberOfFloors {
			checkForButtonPress(orderHandler, i, elevio.BT_HallUp, HALLUP)
			checkForButtonPress(orderHandler, i, elevio.BT_HallDown, HALLDOWN)
			checkForButtonPress(orderHandler, i, elevio.BT_Cab, CAB)
		}
	}
}

