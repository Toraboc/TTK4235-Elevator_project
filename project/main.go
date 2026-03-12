package main

import (
	"fmt"
	"time"

	//. "project/elevator"
	. "project/network"
	. "project/orderHandler"
	. "project/shared"
)

func main() {

	fmt.Println("Starting elevator")
	GetMyId() // Initialize

	orderHandler := NewOrderHandler()

	worldViewMergeChannel := make(chan SyncView)
	connectedNodesUpdateChannel := make(chan NodeIdSet)

	go NetworkProcess(orderHandler, connectedNodesUpdateChannel, worldViewMergeChannel)
	// write two temp goroutines to read from channels and do nothing to prevent blocking
	go func() {
		for range worldViewMergeChannel {
		}
	}()
	go func() {
		for range connectedNodesUpdateChannel {
		}
	}()

	//go ElevatorProcess(orderHandler)

	for {
		time.Sleep(1 * time.Second)
	}
}
