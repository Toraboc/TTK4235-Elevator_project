package orders

import (
	. "project/shared"
)

// Merge our worldview with the incomming SyncMessage
// This will only sync the orders and elevatorStates
func MergeWorldView(SyncMessage SyncMessage) {
	// At last
	// This must also be called if our own elevatorsstate changes
	hallRequestAssigner()
}

func GetWorldview() Worldview {
	return Worldview{}
}

func hallRequestAssigner() {

}

func elevatorStop(floor int) {

}

func getNextTargetFloor() (int, error) {
	return 1, nil
}
