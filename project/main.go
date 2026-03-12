package main

import (
	"fmt"
	"time"

	. "project/elevator"
	. "project/network"
	. "project/orderHandler"
	. "project/shared"
)

func main() {

	fmt.Println("Starting elevator")

	GetMyId() // Initialize

	orderChannels := NewOrderChannels()

	go NetworkProcess(orderChannels)
	go ElevatorProcess(orderChannels)
	go OrderProcess(orderChannels)

	for {
		time.Sleep(1 * time.Second)
	}
}
