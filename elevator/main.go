package main

import (
	"Driver-go/elevio"
	"fmt"
)

func main() {
	fmt.Println("Starting main...")
	elevio.Init("localhost:3333",4)
	Fsm()
}
