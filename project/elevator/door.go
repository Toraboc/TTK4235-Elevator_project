package elevator

import (
	"Driver-go/elevio"
	"fmt"
	"time"
)

type Door struct {
	IsOpen     bool
	closeTime  time.Time
	Obstructed bool
	WillOpen   bool
}

func getDoorState(door *Door) bool {
	return door.IsOpen
}

func openDoor(door *Door) {
	startDoorTimer(door)
	elevio.SetDoorOpenLamp(true)
	door.IsOpen = true
}

func closeDoor(door *Door) {
	elevio.SetDoorOpenLamp(false)
	door.IsOpen = false
}

func startDoorTimer(door *Door) {
	door.closeTime = time.Now().Add(3 * time.Second)
}

func handleObstruction(door *Door) {
	if elevio.GetObstruction() {
		fmt.Println("Door obstructed")
	}
	if elevio.GetObstruction() && door.IsOpen {
		openDoor(door)
	}
}

func DoorModuleLoop(door *Door) {
	handleObstruction(door)
	if door.closeTime.Before(time.Now()) && door.Obstructed == false {
		closeDoor(door)
	}
}
