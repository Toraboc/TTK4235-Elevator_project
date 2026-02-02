package main

type Direction int
const (
    UP Direction = iota
    DOWN
)

type ElevatorBehaviour int
const (
    IDLE ElevatorBehaviour = iota
    MOVING
    PASSENGER_TRANSFER
    OBSTRCTED
)

type ElevatorState struct {
    bahviour  ElevatorBehaviour
    position  int
    direction Direction
}
