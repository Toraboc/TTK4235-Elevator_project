package shared

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
    hallUpOrders [NumberOfFloors]Order
    hallDownOrders [NumberOfFloors]Order
    cabOrders map[NodeId][NumberOfFloors]Order
}

type Worldview struct {
    orders Orders
    elevatorStates map[NodeId]ElevatorState
    assignedHallUpOrders [NumberOfFloors]bool
    assignedHallDownOrders [NumberOfFloors]bool
    assignedCabOrders [NumberOfFloors]bool
}
