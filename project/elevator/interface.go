package elevator

import (
	. "project/orderHandler"
	. "project/shared"
)

type ElevatorInterface struct {
	ElevatorStateCh   chan<- ElevatorState
	OrderCompletedCh  chan<- OrderCompletedEvent
	NewOrderCh        chan<- NewOrderEvent
	ConfirmedOrdersCh <-chan ConfirmedOrders
	TargetFloorCh     <-chan int
}
