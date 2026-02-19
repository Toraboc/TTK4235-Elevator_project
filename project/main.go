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
	GetMyId() // Initialize myId

	go NetworkProcess()

	//ElevatorProcess()

	for {
		time.Sleep(1 * time.Second)
	}
}