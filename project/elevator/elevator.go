package elevator

import . "project/shared"

var elevatorStateDetailed

func ElevatorProcess() {
	positioning := InitPositioning()

	fmt.println("Elevator state is determined.")

	go handleButtons()

	positioning.handleDirving()
}

func GetElevatorState() ElevatorState {
	return ElevatorState{}
}
