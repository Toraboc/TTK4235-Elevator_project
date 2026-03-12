package orderHandler

import (
	"encoding/json"
	"fmt"
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

// run the script and update assigned requests
func (worldView *WorldView) hallRequestAssigner() {

	confirmedOrders := worldView.getConfirmedOrders()

	//Caborders cannot be assigned between elevators
	worldView.AssignedCabOrders = confirmedOrders.Cab

	//exits if no orders are confirmed
	if noOrdersConfirmed(confirmedOrders) {
		worldView.AssignedHallUpOrders = confirmedOrders.HallUp
		worldView.AssignedHallDownOrders = confirmedOrders.HallDown
		return
	}

	hallRequests := make([][2]bool, NumberOfFloors)
	for floor := range NumberOfFloors {
		hallRequests[floor][0] = confirmedOrders.HallUp[floor]
		hallRequests[floor][1] = confirmedOrders.HallDown[floor]
	}

	states := make(map[string]hallRequestAssignerInputState)
	for nodeId := range worldView.ConnectedNodes {
		elevatorState := worldView.ElevatorStates[nodeId]

		// Skip elevators that cannot be assigned Orders
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
		return
	}

	input := hallRequestAssignerInput{
		HallRequests: hallRequests,
		States:       states,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		panic("hallRequestAssigner: failed to marshal input: " + err.Error())
	}

	fmt.Println(string(inputJSON))

	command := exec.Command("./hall_request_assigner", "--input", string(inputJSON))
	outputJSON, err := command.Output()
	if err != nil {
		panic("hallRequestAssigner: command failed: " + err.Error())
	}

	var hallAssignmentsByElevator map[string][][]bool
	if err := json.Unmarshal(outputJSON, &hallAssignmentsByElevator); err != nil {
		panic("hallRequestAssigner: failed to unmarshal output: " + err.Error())
	}

	myId := GetMyId()
	assignedHallRequests, exists := hallAssignmentsByElevator[myId.String()]
	if !exists {
		return
	}

	worldView.AssignedHallUpOrders = [NumberOfFloors]bool{}
	worldView.AssignedHallDownOrders = [NumberOfFloors]bool{}

	for floor := range NumberOfFloors {
		worldView.AssignedHallUpOrders[floor] = assignedHallRequests[floor][0]
		worldView.AssignedHallDownOrders[floor] = assignedHallRequests[floor][1]

	}
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
