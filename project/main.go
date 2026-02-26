package main

import (
	"fmt"
	"time"

	//. "project/elevator"
	. "project/network"
	. "project/shared"
)

func main() {

	fmt.Println("Starting elevator")
	GetMyId() // Initialize 
	
	orderHandler := NewOrderHandler()

	go NetworkProcess(orderHandler)

	orderHandler.GetWorldView()

	//ElevatorProcess(orderHandler)

	for {
		time.Sleep(1 * time.Second)
	}
}