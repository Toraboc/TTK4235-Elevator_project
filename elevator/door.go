package main

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

type Door struct {
	isOpen       bool
	closeTime    time.Time
	Obstructed   bool
	willOpen bool
}

func getDoorState(door *Door) bool {
	return door.isOpen
}

func openDoor(door *Door) {
	startDoorTimer(door)
	elevio.SetDoorOpenLamp(true)
	door.isOpen = true
}

func closeDoor(door *Door) {
	elevio.SetDoorOpenLamp(false)
	door.isOpen = false
}

func startDoorTimer(door *Door) {
	door.closeTime = time.Now().Add(3 * time.Second)
}

func handleObstruction(door *Door) {
	if elevio.GetObstruction() {
		fmt.Println("Door obstructed")
	}
	if elevio.GetObstruction() && door.isOpen {
		openDoor(door)
	}
}

func doorModuleLoop(door *Door) {
	handleObstruction(door)
	if door.closeTime.Before(time.Now()) && door.Obstructed == false {
		closeDoor(door)
	}
}
