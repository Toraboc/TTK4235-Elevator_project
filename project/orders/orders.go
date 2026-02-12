package orders

import (
	. "project/shared"
)

var WorldView Worldview

// Merge our worldview with the incomming data in some way, not dependent on network
// This will only sync the orders and elevatorStates
func MergeWorldView(SyncMessage SyncMessage) {

	// At last
	// This must also be called if our own elevatorsstate changes
	hallRequestAssigner()
}

func GetWorldview() Worldview {
	return WorldView
}

func UpdateConnectedNodes(ids []NodeId) {
	WorldView.ConnectedNodes = ids
}

func hallRequestAssigner() {

}

// This function will receive updates from the elevator
func ElevatorStateChange(state ElevatorState) {

}

// The datainout here will we figure out later
func newOrder() {

}

// Return the next target floor
func GetNextTargetFloor() (int, error) {
	return 1, nil
}
