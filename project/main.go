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

	targetFloorCh := make(chan int, 1)
	elevatorStateCh := make(chan ElevatorState, 1)
	orderCompletedCh := make(chan OrderCompleted, 10)
	orderNewCh := make(chan OrderNew, 10)
	worldViewMergeCh := make(chan SyncView, 1)
	connectedNodesUpdateCh := make(chan NodeIdSet, 1)

	// write two temp goroutines to read from channels and do nothing to prevent blocking
	go func() {
		for range worldViewMergeCh {
		}
	}()
	go func() {
		for range connectedNodesUpdateCh {
		}
	}()

	orderHandler := NewOrderHandler(targetFloorCh, elevatorStateCh, orderCompletedCh, orderNewCh)

	go NetworkProcess(orderHandler, connectedNodesUpdateCh, worldViewMergeCh)

	go ElevatorProcess(*elevatorServerHost, orderHandler, elevatorStateCh, orderCompletedCh, targetFloorCh, orderNewCh)

	for {
		time.Sleep(1 * time.Second)
	}
}
