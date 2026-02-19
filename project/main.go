package main

import (
	"fmt"
	// "time"

	. "project/elevator"
	// . "project/network"
)

func main() {

	fmt.Println("Starting elevator")

	// go NetworkProcess()

	ElevatorProcess()

	// for {
	// 	time.Sleep(1 * time.Second)
	// }
}
