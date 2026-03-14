package main

import (
	"flag"
	"fmt"
	"time"

	. "project/elevator"
	. "project/network"
	. "project/orderHandler"
	. "project/shared"
)

func main() {

	elevatorServerHost := flag.String("server", "localhost:15657", "Elevator server host")
	flag.Parse()

	fmt.Println("Starting elevator")

	InitMyId()

	orderChannels := NewOrderChannels()

	go NetworkProcess(orderChannels)
	go ElevatorProcess(orderChannels, *elevatorServerHost)
	go OrderProcess(orderChannels)

	for {
		time.Sleep(1 * time.Second)
	}
}
