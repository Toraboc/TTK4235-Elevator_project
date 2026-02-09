package main

// Merge our worldview with the incomming SyncMessage
// This will only sync the orders and elevatorStates
func mergeWorldView(SyncMessage SyncMessage) {

	// At last
	// This must also be called if our own elevatorsstate changes
	hallRequestAssigner()
}

func hallRequestAssigner() {

}

func elevatorStop(floor int) {

}

func getNextTargetFloor() (int, error) {
	return 1, nil
}
