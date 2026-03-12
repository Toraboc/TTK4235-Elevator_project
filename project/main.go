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

	targetFloorCh := make(chan int, 1)
	elevatorStateCh := make(chan ElevatorState, 1)
	orderCompletedCh := make(chan OrderCompleted, 10)
	orderNewCh := make(chan OrderNew, 10)
	worldViewMergeChannel := make(chan SyncView, 1)
	connectedNodesUpdateChannel := make(chan NodeIdSet, 1)

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
