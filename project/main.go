package main

import (
	"fmt"

	//. "project/elevator"
	. "project/network"
)

func main() {

	fmt.Println("Starting elevator")

	fmt.Printf("My Ip: %s\n", NodeIdtoString(GetOwnId()))
	//go NetworkProcess()

	//ElevatorProcess()

}
