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
	worldViewMergeCh := make(chan SyncView)
	connectedNodesUpdateCh := make(chan NodeIdSet)

	// write two temp goroutines to read from channels and do nothing to prevent blocking
	go func() {
		for range worldViewMergeCh {
		}
	}()
	go func() {
		for range connectedNodesUpdateCh {
		}
	}()

	orderHandler := NewOrderHandler(targetFloorCh, elevatorStateCh, orderCompletedCh)

	go NetworkProcess(orderHandler, connectedNodesUpdateCh, worldViewMergeCh)

	go ElevatorProcess(orderHandler, elevatorStateCh, orderCompletedCh, targetFloorCh)

	for {
		time.Sleep(1 * time.Second)
	}
}
