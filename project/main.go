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
	udpPort := flag.Int("port", 44039, "The port to listen and send UDP packages")
	myIdSuffix := flag.Int("id", 0, "The last octet in the nodeID. The id will be prefixed with 10.100.23.")
	flag.Parse()

	fmt.Println("Starting elevator")

	// Setup the nodeId
	if *myIdSuffix != 0 {
		myId := (NodeId)((10 << 24) | (100 << 16) | (23 << 8) | (*myIdSuffix & 0xff))
		SetMyId(myId)
	} else {
		GetMyId()
	}

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

	go NetworkProcess(*udpPort, orderHandler, connectedNodesUpdateChannel, worldViewMergeChannel)

	go ElevatorProcess(*elevatorServerHost, orderHandler, elevatorStateCh, orderCompletedCh, targetFloorCh, orderNewCh)

	for {
		time.Sleep(1 * time.Second)
	}
}
