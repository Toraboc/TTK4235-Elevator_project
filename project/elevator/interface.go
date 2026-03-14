package elevator

import (
	. "project/shared"
	. "project/orderHandler"
)

type ElevatorInterface struct {
	ElevatorStateCh        chan<- ElevatorState
	OrderCompletedCh       chan<- OrderCompletedEvent
	NewOrderCh             chan<- NewOrderEvent
	ConfirmedOrdersCh      <-chan ConfirmedOrders
	TargetFloorCh          <-chan int
}
