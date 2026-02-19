package orders

import (
	"encoding/json"
	"fmt"
	"os"
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
	HallRequests [][2]bool                        `json:"hallRequests"`
	States       map[string]hallRequestAssignerInputState `json:"states"`
}

//run the script and update assigned requests
func (worldView *WorldView) hallRequestAssigner() {

	confirmedOrders := worldView.GetConfirmedOrders()

	worldView.AssignedCabOrders = confirmedOrders.Cab

	hallRequests := make([][2]bool, NumberOfFloors)
	for floor := 0; floor < NumberOfFloors; floor++ {
		hallRequests[floor][0] = confirmedOrders.HallUp[floor]
		hallRequests[floor][1] = confirmedOrders.HallDown[floor]
	}

	states := make(map[string]hallRequestAssignerInputState)
	for nodeId := range worldView.ConnectedNodes {
		cabRequests := make([]bool, NumberOfFloors)
		if cabOrderList, exists := worldView.Orders.CabOrders[nodeId]; exists {
			cabRequests = findConfirmedOrdersInArray(cabOrderList, nodeId)[:]
		}

		states[nodeId.String()] = hallRequestAssignerInputState{
			Behaviour:   "idle",
			Floor:       0,
			Direction:   "stop",
			CabRequests: cabRequests,
		}
	}

	input := hallRequestAssignerInput{
		HallRequests: hallRequests,
		States:       states,
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		fmt.Println("hallRequestAssigner: failed to marshal input:", err)
		return
	}

	command := exec.Command("./hall_request_assigner", "--input", string(inputJSON))
	outputJSON, err := command.Output()
	if err != nil {
		fmt.Println("hallRequestAssigner: command failed:", err)
		return
	}

	var hallAssignmentsByElevator map[string][][]bool
	if err := json.Unmarshal(outputJSON, &hallAssignmentsByElevator); err != nil {
		fmt.Println("hallRequestAssigner: failed to unmarshal output:", err)
		return
	}

	ownId := getOwnId()
	assignedHallRequests, exists := hallAssignmentsByElevator[ownId.String()]
	if !exists {
		return
	}

	worldView.AssignedHallUpOrders = [NumberOfFloors]bool{}
	worldView.AssignedHallDownOrders = [NumberOfFloors]bool{}

	for floor := 0; floor < NumberOfFloors; floor++ {
		worldView.AssignedHallUpOrders[floor] = assignedHallRequests[floor][0]
		worldView.AssignedHallDownOrders[floor] = assignedHallRequests[floor][1]
	
	}
}

