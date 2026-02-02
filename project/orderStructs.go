package main

import (
	"time"
)

type OrderStatus int
const (
    NEW OrderStatus = iota
    COMPLETED
)

type Order struct {
    lastEvent OrderStatus // skal dette vere OrderStatus?
    lastUpdate time.Time
    confirmedBy []NodeId
}
 
type Orders struct {
    hallUpOrders [numberOfFloors]Order
    hallDownOrders [numberOfFloors]Order
    cabOrders map[NodeId][numberOfFloors]Order
}

type Worldview struct {
    orders Orders
    elevatorStates map[NodeId]ElevatorState
    assignedHallUpOrders [numberOfFloors]bool
    assignedHallDownOrders [numberOfFloors]bool
    assignedCabOrders [numberOfFloors]bool
}


