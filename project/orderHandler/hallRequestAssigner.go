package orderHandler

import (
	"encoding/json"
	"os/exec"
	. "project/shared"
)

type hallRequestAssignerInputState struct {
	Behaviour   string               `json:"behaviour"`
	Floor       int                  `json:"floor"`
	Direction   string               `json:"direction"`
	CabRequests [NumberOfFloors]bool `json:"cabRequests"`
}

type hallRequestAssignerInput struct {
	HallRequests [NumberOfFloors][2]bool                  `json:"hallRequests"`
	States       map[string]hallRequestAssignerInputState `json:"states"`
}

type AssignedOrders struct {
	HallUp   [NumberOfFloors]bool
	HallDown [NumberOfFloors]bool
	Cab      [NumberOfFloors]bool
}

func hallRequestAssigner(worldView *WorldView, nodeId NodeId) AssignedOrders {
	var assignedOrders AssignedOrders

	confirmedOrders := getConfirmedOrders(worldView.Orders[nodeId], nodeId)

	if !anyOrdersConfirmed(confirmedOrders) {
		return assignedOrders
	}

	assignedOrders.Cab = confirmedOrders.Cab

	var hallRequests [NumberOfFloors][2]bool
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

		var cabRequests [NumberOfFloors]bool
		if orders, exists := worldView.Orders[nodeId]; exists {
			if cabOrders, exists := orders.CabOrders[nodeId]; exists {
				cabRequests = findConfirmedOrdersInArray(cabOrders)
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
		panic("Failed to marshal input for the Hall Request Assigner: " + err.Error())
	}

	command := exec.Command("./hall_request_assigner", "--input", string(inputJSON))
	assignerOutput, err := command.CombinedOutput()
	if err != nil {
		panic("hall_request_assigner: command failed:\n" + err.Error() + "\n\nInput:\n" + string(inputJSON) + "\n\nOutput:\n" + string(assignerOutput))
	}

	var hallAssignmentsByElevator map[string][][]bool
	if err := json.Unmarshal(assignerOutput, &hallAssignmentsByElevator); err != nil {
		panic("Failed to unmarshal output from the Hall Request Assigner: " + err.Error())
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
