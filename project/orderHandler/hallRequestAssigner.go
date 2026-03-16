package orderHandler

import (
	"encoding/json"
	"os/exec"
	. "project/shared"
)

type hallRequestAssignerInputState struct {
	Behaviour   string `json:"behaviour"`
	Floor       int    `json:"floor"`
	Direction   string `json:"direction"`
	CabRequests []bool `json:"cabRequests"`
}

type hallRequestAssignerInput struct {
	HallRequests [][2]bool                                `json:"hallRequests"`
	States       map[string]hallRequestAssignerInputState `json:"states"`
}

type AssignedOrders struct {
	HallUp   [NumberOfFloors]bool
	HallDown [NumberOfFloors]bool
	Cab      [NumberOfFloors]bool
}

func hallRequestAssigner(worldView *WorldView, nodeId NodeId) AssignedOrders {
	var assignedOrders AssignedOrders

	confirmedOrders := worldView.getConfirmedOrders()

	assignedOrders.Cab = confirmedOrders.Cab

	if !anyOrdersConfirmed(confirmedOrders) {
		return assignedOrders
	}

	hallRequests := make([][2]bool, NumberOfFloors)
	for floor := range NumberOfFloors {
		hallRequests[floor][0] = confirmedOrders.HallUp[floor]
		hallRequests[floor][1] = confirmedOrders.HallDown[floor]
	}

	states := make(map[string]hallRequestAssignerInputState)
	for nodeId := range worldView.ConnectedNodes {
		elevatorState := worldView.ElevatorStates[nodeId]

		if !elevatorState.Behaviour.CanBeAssignedOrders() {
			continue
		}

		cabRequests := make([]bool, NumberOfFloors)
		if orders, exists := worldView.Orders[nodeId]; exists {
			if cabOrders, exists := orders.CabOrders[nodeId]; exists {
				confirmedCab := findConfirmedOrdersInArray(cabOrders)
				cabRequests = confirmedCab[:]
			}
		}

		states[nodeId.String()] = hallRequestAssignerInputState{
			Behaviour:   getBehaviourString(elevatorState),
			Floor:       elevatorState.Floor,
			Direction:   getDirectionString(elevatorState),
			CabRequests: cabRequests,
		}
	}

	if len(states) == 0 {
		return assignedOrders
	}

	input := hallRequestAssignerInput{
		HallRequests: hallRequests,
		States:       states,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		panic("hallRequestAssigner: failed to marshal input: " + err.Error())
	}

	command := exec.Command("./hall_request_assigner", "--input", string(inputJSON))
	outputJSON, err := command.Output()
	if err != nil {
		panic("hallRequestAssigner: command failed: " + err.Error())
	}

	var hallAssignmentsByElevator map[string][][]bool
	if err := json.Unmarshal(outputJSON, &hallAssignmentsByElevator); err != nil {
		panic("hallRequestAssigner: failed to unmarshal output: " + err.Error())
	}

	assignedHallRequests, exists := hallAssignmentsByElevator[nodeId.String()]
	if !exists {
		return assignedOrders
	}

	for floor := range NumberOfFloors {
		assignedOrders.HallUp[floor] = assignedHallRequests[floor][0]
		assignedOrders.HallDown[floor] = assignedHallRequests[floor][1]

	}

	return assignedOrders
}

func getBehaviourString(elevatorState ElevatorState) string {
	switch elevatorState.Behaviour {
	case MOVING:
		return "moving"
	case IDLE:
		return "idle"
	case PASSENGER_TRANSFER:
		return "doorOpen"
	default:
		panic("Cannot determine behaviour for an elevator that is not either MOVING, IDLE or PASSENGER_TRANSFER")
	}
}

func getDirectionString(elevatorState ElevatorState) string {
	if !elevatorState.Behaviour.CanBeAssignedOrders() {
		panic("Cannot determine direction for an elevator that is not working properly.")
	}
	if elevatorState.Behaviour == MOVING {
		switch elevatorState.Direction {
		case UP:
			return "up"
		case DOWN:
			return "down"
		default:
			panic("Invalid elevator direction.")
		}
	}

	return "stop"
}
