package orders

import (
	. "project/shared"
)

type OrderStatus int
const (
    NO_ORDER OrderStatus = iota
    UNCONFIRMED
    CONFIRMED
    FINISHED
)

type OrderList [NumberOfFloors]OrderStatus

type Orders struct {
	HallUpOrders   OrderList
	HallDownOrders OrderList
	CabOrders      map[NodeId]OrderList
}

type WorldView struct {
	Orders                 map[NodeId]Orders
	ConnectedNodes         NodeIdSet
	ElevatorStates         map[NodeId]ElevatorState
	AssignedHallUpOrders   [NumberOfFloors]bool
	AssignedHallDownOrders [NumberOfFloors]bool
	AssignedCabOrders      [NumberOfFloors]bool
}
