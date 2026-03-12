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

	targetFloorCh := make(chan int)
	elevatorStateCh := make(chan ElevatorState)
	orderCompletedCh := make(chan OrderCompleted)
	orderNewCh := make(chan OrderNew)
	worldViewMergeChannel := make(chan SyncView)
	connectedNodesUpdateChannel := make(chan NodeIdSet)

	// write two temp goroutines to read from channels and do nothing to prevent blocking
	go func() {
		for range worldViewMergeChannel {
		}
	}()
	go func() {
		for range connectedNodesUpdateChannel {
		}
	}()

	orderHandler := NewOrderHandler(targetFloorCh, elevatorStateCh, orderCompletedCh, orderNewCh)

	go NetworkProcess(orderHandler, connectedNodesUpdateChannel, worldViewMergeChannel)

	go ElevatorProcess(orderHandler, elevatorStateCh, orderCompletedCh, targetFloorCh, orderNewCh)

	for {
		time.Sleep(1 * time.Second)
	}
}
