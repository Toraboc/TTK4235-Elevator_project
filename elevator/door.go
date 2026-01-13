package main

import (
	"Driver-go/elevio"
	"time"
	"fmt"
)

type DoorState int

const (
	DoorOpen DoorState = iota
	DoorClosed
)

type Door struct {
	isOpen DoorState
	closeTime time.Time
	Obstructed bool
	willOpenDoor bool
}




func getDoorState(door *Door) DoorState {
	return door.isOpen
}


func openDoor(door *Door) {
	startDoorTimer(door)
	elevio.SetDoorOpenLamp(true)
	door.isOpen = DoorOpen
}


func closeDoor(door *Door) {
	elevio.SetDoorOpenLamp(false)
	door.isOpen = DoorClosed
}


func startDoorTimer(door *Door) {
	door.closeTime = time.Now().Add(3 * time.Second)
}


func handleObstruction(door *Door) {
	if elevio.GetObstruction() {
		fmt.Println("Door obstructed")
	}
	if elevio.GetObstruction() && door.isOpen == DoorOpen {
		elevio.SetDoorOpenLamp(true)
	}
}


func DoorModuleLoop(door *Door) {
	handleObstruction(door)
	if door.closeTime.Before(time.Now()) && door.Obstructed == false {
		closeDoor(door)
	}
}