package orders

import (
    "time"
    . "project/shared"
)

type OrderStatus int
const (
    NO_ORDER OrderStatus = iota
    UNCONFIRMED
    CONFIRMED
    FINISHED
)

type Order map[NodeId]OrderStatus
 
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
