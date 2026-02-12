package orders

import (
    "time"
)

type OrderStatus int
const (
    NEW OrderStatus = iota
    COMPLETED
)

type Order struct {
    LastEvent OrderStatus // skal dette vere OrderStatus?
    LastUpdate time.Time
    ConfirmedBy NodeIdSet
}
 
type Orders struct {
    HallUpOrders [NumberOfFloors]Order
    HallDownOrders [NumberOfFloors]Order
    CabOrders map[NodeId][NumberOfFloors]Order
}

type WorldView struct {
    Orders Orders
    ConnectedNodes NodeIdSet
    ElevatorStates map[NodeId]ElevatorState
    AssignedHallUpOrders [NumberOfFloors]bool
    AssignedHallDownOrders [NumberOfFloors]bool
    AssignedCabOrders [NumberOfFloors]bool
}
