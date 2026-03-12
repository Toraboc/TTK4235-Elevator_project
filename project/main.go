package main

import (
	"fmt"
	"time"

	. "project/elevator"
	// . "project/network"
	. "project/shared"
	. "project/orderHandler"
)

// func targetFloors(trgFChr chan<- int) {
// 	time.Sleep(1 * time.Second)
//
// 	fmt.Println("New target: 1")
// 	trgFChr <- 1
//
// 	time.Sleep(10 * time.Second)
//
// 	fmt.Println("New target: 3")
// 	trgFChr <- 3
//
// 	time.Sleep(10 * time.Second)
//
// 	fmt.Println("New target: 0")
// 	trgFChr <- 0
// }

func main() {

	fmt.Println("Starting elevator")
	GetMyId() // Initialize 

	targetFloorCh := make(chan int)
	elevatorStateCh := make(chan ElevatorState)
	orderCompletedCh := make(chan OrderCompleted)
	
	orderHandler := NewOrderHandler(targetFloorCh, elevatorStateCh, orderCompletedCh)

	// go NetworkProcess(orderHandler)

	go ElevatorProcess(orderHandler, elevatorStateCh, orderCompletedCh, targetFloorCh)

	for {
		time.Sleep(1 * time.Second)
	}
}
