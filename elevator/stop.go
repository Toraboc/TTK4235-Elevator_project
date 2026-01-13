package main

import (
	"Driver-go/elevio"

)

func activateStopLamp() {
	elevio.SetStopLamp(true)
}

func deactivateStopLamp() {
	elevio.SetStopLamp(false)
}

func handleStop(pos *Position, door *Door, orderHandler *OrderHandler) {
	// Go through emergency stop procedure
	elevio.SetMotorDirection(elevio.MD_Stop)
	activateStopLamp()
	clearAllOrders(orderHandler)
	if elevio.GetFloor() != -1 {
		openDoor(door)
	}
	for elevio.GetStop() {} // Should probably make a fix here later
	deactivateStopLamp()
	startDoorTimer(door)
	pos.targetFloor = -1
}

func stopModuleLoop(pos *Position, door *Door, orderHandler *OrderHandler) {
	if elevio.GetStop() {
		handleStop(pos, door, orderHandler)
	}
}