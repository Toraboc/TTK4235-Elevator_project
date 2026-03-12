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

	orderHandler := NewOrderHandler()

	WorldViewMergeChannel := make(chan SyncView)
	ConnectedNodesUpdateChannel := make(chan NodeIdSet)

	go NetworkProcess(orderHandler, ConnectedNodesUpdateChannel, WorldViewMergeChannel)

	go ElevatorProcess(orderHandler)

	for {
		time.Sleep(1 * time.Second)
	}
}
