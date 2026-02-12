package main

import (
	"fmt"

	. "project/elevator"
	. "project/network"
)

func main() {

	fmt.Println("Starting elevator")

	go NetworkProcess()

	ElevatorProcess()

}
