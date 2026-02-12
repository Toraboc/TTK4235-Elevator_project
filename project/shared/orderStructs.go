package shared

import (
	"sync"
	"time"
)

type OrderStatus int

const (
	NEW OrderStatus = iota
	COMPLETED
)

type Order struct {
	LastEvent   OrderStatus // skal dette vere OrderStatus?
	LastUpdate  time.Time
	ConfirmedBy []NodeId
}

type Orders struct {
	HallUpOrders   [NumberOfFloors]Order
	HallDownOrders [NumberOfFloors]Order
	CabOrders      map[NodeId][NumberOfFloors]Order
}

type Worldview struct {
	Mu                     sync.Mutex
	Orders                 Orders
	ConnectedNodes         []NodeId
	ElevatorStates         map[NodeId]ElevatorState
	AssignedHallUpOrders   [NumberOfFloors]bool
	AssignedHallDownOrders [NumberOfFloors]bool
	AssignedCabOrders      [NumberOfFloors]bool
}
