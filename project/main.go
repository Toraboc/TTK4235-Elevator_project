package main

import (
	"fmt"
	"time"

	. "project/elevator"
	// . "project/network"
	. "project/shared"
	. "project/orderHandler"
)

func targetFloors(trgFChr chan<- int) {
	time.Sleep(1 * time.Second)

	fmt.Println("New target: 1")
	trgFChr <- 1

	time.Sleep(10 * time.Second)

	fmt.Println("New target: 3")
	trgFChr <- 3

	time.Sleep(10 * time.Second)

	fmt.Println("New target: 0")
	trgFChr <- 0
}

func main() {

	fmt.Println("Starting elevator")
	GetMyId() // Initialize 
	
	orderHandler := NewOrderHandler()

	// go NetworkProcess(orderHandler)

	targetFloorCh := make(chan int)

	go ElevatorProcess(orderHandler, targetFloorCh)

	targetFloors(targetFloorCh)

	for {
		time.Sleep(1 * time.Second)
	}
}
