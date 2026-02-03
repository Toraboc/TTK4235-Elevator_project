package main

import (
	"Driver-go/elevio"
	"fmt"
	"project/elevator"
	"sync"
	//"project/drivers/network"
	//"project/drivers/orderhandler"
)

var worldview Worldview
var worldviewMutex sync.Mutex
var knowsMe KnowsMe

func main() {
	//elevio.Init("localhost:15657", 4)

	fmt.Println("Starting elevator")

	go networkProcess()

	//elevatorProcess()

	// Heiskode greier under ellerno
	for {
	}
}

func mainOld() {
	elevio.Init("localhost:15657", 4)

	fmt.Println("Starting elevator")

	orderHandler := elevator.OrderModuleInit()
	door := elevator.Door{IsOpen: false, Obstructed: false, WillOpen: false}

	position := elevator.Position{}
	elevator.InitPosition(&position)

	//door.GetDoorState()

	for {
		elevator.DoorModuleLoop(&door)
		elevator.StopModuleLoop(&position, &door, &orderHandler)
		elevator.OrderModuleLoop(&orderHandler)
		elevator.PositionModuleLoop(&position, &door)

		nextOrderFloor := elevator.GetNextOrder(&orderHandler, position.LastFloor, position.LastDirection)
		if elevator.GetDirection(&position) == elevator.DirStop && door.IsOpen {
			elevator.StoppedAtFloor(&orderHandler, position.LastFloor, nextOrderFloor)
		}

		if nextOrderFloor != -1 && nextOrderFloor != position.TargetFloor {
			fmt.Printf("Going to floor: %d\n", nextOrderFloor)
			elevator.GotoFloor(&position, &door, nextOrderFloor)
		}
	}
}
