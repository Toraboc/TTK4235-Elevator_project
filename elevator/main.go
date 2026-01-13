package main

import (
	"Driver-go/elevio"
	"fmt"
)

var N_FLOORS = 4

func main() {
	elevio.Init("localhost:3333",4)

	fmt.Println("Starting elevator")

	orderHandler := orderModuleInit()
	door := Door{isOpen: DoorClosed, Obstructed: false, willOpenDoor: false}
	position := Position{}
	initPosition(&position)

	for {
		DoorModuleLoop(&door)
		stopModuleLoop(&position, &door, &orderHandler)
		orderModuleLoop(&orderHandler)
		positionModuleLoop(&position, &door)

		nextOrderFloor := getNextOrder(&orderHandler, position.lastFloor, position.lastDirection)
		if getDirection(&position) == DirStop && door.isOpen == DoorOpen {
			stoppedAtFloor(&orderHandler, position.lastFloor)
		}

		if nextOrderFloor != -1 && nextOrderFloor != position.targetFloor {
			fmt.Println("Going to floor: %d\n", nextOrderFloor)
			gotoFloor(&position, &door, nextOrderFloor)
		}
	}
}
