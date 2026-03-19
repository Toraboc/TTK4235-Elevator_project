package orderHandler

import (
	"fmt"
	. "project/shared"
)

func getNextTargetFloor(worldView WorldView, nodeId NodeId) (int, error) {
	elevatorState, exists := worldView.ElevatorStates[nodeId]
	if !exists {
		return -1, fmt.Errorf("missing elevatorState for own node")
	}

	floor := elevatorState.Floor
	if floor < 0 || floor >= NumberOfFloors {
		return -1, fmt.Errorf("invalid current floor: %d", floor)
	}

	orders := hallRequestAssigner(worldView, nodeId)

	if elevatorState.Direction == UP {
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if orders.Cab[floor] || orders.HallUp[floor] {
				return floor, nil
			}
		}

		for i := floor + 1; i < NumberOfFloors; i++ {
			if orders.Cab[i] || orders.HallUp[i] {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= 0; i-- {
			if orders.Cab[i] || orders.HallDown[i] {
				return i, nil
			}
		}
		for i := 0; i <= floor; i++ {
			if orders.Cab[i] || orders.HallUp[i] {
				return i, nil
			}
		}
	}
	if elevatorState.Direction == DOWN {
		if elevatorState.Behaviour == PASSENGER_TRANSFER || elevatorState.Behaviour == IDLE {
			if orders.Cab[floor] || orders.HallDown[floor] {
				return floor, nil
			}
		}

		for i := floor - 1; i >= 0; i-- {
			if orders.Cab[i] || orders.HallDown[i] {
				return i, nil
			}
		}
		for i := range NumberOfFloors {
			if orders.Cab[i] || orders.HallUp[i] {
				return i, nil
			}
		}
		for i := NumberOfFloors - 1; i >= floor; i-- {
			if orders.Cab[i] || orders.HallDown[i] {
				return i, nil
			}
		}
	}
	return -1, nil
}
