package elevator

import (
	"Driver-go/elevio"
	. "project/shared"
)

func (lightStatus *OrderButtons) updateLights(currentOrders OrderButtons) {
	updateSingleLight := func (floor int, matrixIndex int, lamp elevio.ButtonType) {
		if lightStatus[floor][matrixIndex] != currentOrders[floor][matrixIndex] {
			lightStatus[floor][matrixIndex] = currentOrders[floor][matrixIndex]
			elevio.SetButtonLamp(lamp, floor, lightStatus[floor][matrixIndex])
		}
	}
	for i := range NumberOfFloors {
		updateSingleLight(i, 0, elevio.BT_HallUp)
		updateSingleLight(i, 1, elevio.BT_HallDown)
		updateSingleLight(i, 2, elevio.BT_Cab)
	}
}

var lightStatus OrderButtons

func ordersToButtonMatrix(hallUpOrders [NumberOfFloors]Order, hallDownOrders [NumberOfFloors]Order, cabOrders [NumberOfFloors]Order) OrderButtons {
	var buttonMatrix OrderButtons
	for i := range NumberOfFloors {
		// BUG: This is wrong, the orders needs to be confirmed by everyone.
		buttonMatrix[i][0] = hallUpOrders[i].LastEvent == NEW
		buttonMatrix[i][1] = hallDownOrders[i].LastEvent == NEW
		buttonMatrix[i][2] = cabOrders[i].LastEvent == NEW
	}

	return buttonMatrix
}

func SetOrderLights(hallUpOrders [NumberOfFloors]Order, hallDownOrders [NumberOfFloors]Order, cabOrders [NumberOfFloors]Order) {
	lightStatus.updateLights(ordersToButtonMatrix(hallUpOrders, hallDownOrders, cabOrders))
}

