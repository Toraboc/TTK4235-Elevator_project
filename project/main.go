package main

import (
	"fmt"
	"time"

	. "project/elevator"
	. "project/network"
	. "project/shared"
	. "project/orderHandler"
)

func main() {

	fmt.Println("Starting elevator")
	GetMyId() // Initialize 
	
	orderHandler := NewOrderHandler()

	go NetworkProcess(orderHandler)

	go ElevatorProcess(orderHandler)

	for {
		time.Sleep(1 * time.Second)
	}
}
