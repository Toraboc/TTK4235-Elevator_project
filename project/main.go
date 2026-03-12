package main

import (
	"fmt"
	"time"
	"flag"

	. "project/elevator"
	. "project/network"
	. "project/orderHandler"
	. "project/shared"
)

func main() {

	elevatorServerHost := flag.String("server", "localhost:15657", "Elevator server host")
	flag.Parse()

	fmt.Println("Starting elevator")

	GetMyId() // Setup the nodeId

	orderChannels := NewOrderChannels()

	go NetworkProcess(orderChannels)
	go ElevatorProcess(*elevatorServerHost, orderChannels)
	go OrderProcess(orderChannels)

	for {
		time.Sleep(1 * time.Second)
	}
}
