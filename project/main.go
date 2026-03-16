package main

import (
	"flag"
	"fmt"
	"time"

	. "project/elevator"
	. "project/network"
	. "project/orderHandler"
	. "project/shared"
)

func main() {

	elevatorServerHost := flag.String("server", "localhost:15657", "Elevator server host")
	flag.Parse()

	fmt.Println("Starting elevator")

	InitMyId()

	ConnectedNodesUpdateCh := make(chan NodeIdSet, 1)
	SyncMergeCh := make(chan SyncData, 1)
	ElevatorStateCh := make(chan ElevatorState, 1)
	OrderCompletedCh := make(chan OrderCompletedEvent, 10)
	NewOrderCh := make(chan NewOrderEvent, 10)
	RequestSyncCh := make(chan chan SyncData)
	ConfirmedOrdersCh := make(chan ConfirmedOrders, 1)
	TargetFloorCh := make(chan int, 1)

	orderHandlerChannels := OrderHandlerInterface{
		ConnectedNodesUpdateCh: ConnectedNodesUpdateCh,
		SyncMergeCh:            SyncMergeCh,
		ElevatorStateCh:        ElevatorStateCh,
		OrderCompletedCh:       OrderCompletedCh,
		NewOrderCh:             NewOrderCh,
		ConfirmedOrdersCh:      ConfirmedOrdersCh,
		TargetFloorCh:          TargetFloorCh,
		RequestSyncCh:          RequestSyncCh,
	}

	elevatorChannels := ElevatorInterface{
		ElevatorStateCh:   ElevatorStateCh,
		OrderCompletedCh:  OrderCompletedCh,
		NewOrderCh:        NewOrderCh,
		ConfirmedOrdersCh: ConfirmedOrdersCh,
		TargetFloorCh:     TargetFloorCh,
	}

	networkChannels := NetworkInterface{
		ConnectedNodesUpdateCh: ConnectedNodesUpdateCh,
		SyncMergeCh:            SyncMergeCh,
		RequestSyncCh:          RequestSyncCh,
	}

	go NetworkProcess(networkChannels)
	go ElevatorProcess(elevatorChannels, *elevatorServerHost)
	go OrderProcess(orderHandlerChannels)

	for {
		time.Sleep(1 * time.Second)
	}
}
