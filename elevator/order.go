package main

import (
	"Driver-go/elevio"
)

type OrderHandler struct {
	floorState [4] int // [floor][buttonType] ------THEODOR: WE NEED TO FIX ikke lagre som bits!!!!!!!-----------------
}

func orderModuleInit() OrderHandler {
	var ret OrderHandler
	clearAllOrders(&ret)
	return ret
}

func convertBtnTypeToDirection(btnType elevio.ButtonType) Direction {
	switch btnType {
		case elevio.BT_Cab:
			return DirStop
		case elevio.BT_HallUp:
			return DirUp
		case elevio.BT_HallDown:
			return DirDown
		default:
			return DirStop
	}
}

func convertDirectionToBtnType(dir Direction) elevio.ButtonType {
	switch dir {
		case DirUp:
			return elevio.BT_HallUp
		case DirDown:
			return elevio.BT_HallDown
		case DirStop:
			return elevio.BT_Cab
		default:
			return elevio.BT_Cab
	}
}

func directionToBitMap(dir Direction) int { // Theodor explain yourself------------------------------------------
	switch dir {
		case DirStop:
			return 1
		case DirUp:
			return 2
		case DirDown:
			return 4
		default:
			return 1
	}
}



func addOrder(orderHandler *OrderHandler, floor int, dir Direction) {
	newState := orderHandler.floorState[floor] | directionToBitMap(dir)
	if newState == orderHandler.floorState[floor] {
		return
	}

	orderHandler.floorState[floor] = newState

	elevio.SetButtonLamp(convertDirectionToBtnType(dir), floor, true)
}

func orderModuleLoop(orderHandler *OrderHandler) {
	for floor := range N_FLOORS {
		for button := range 3 {
			if elevio.GetButton(elevio.ButtonType(button), floor) {
				addOrder(orderHandler, floor, convertBtnTypeToDirection(elevio.ButtonType(button)))
			}
		}
	}
}

func getNextOrder(orderHandler *OrderHandler, lastFloor int, drivingDirection Direction) int {
	
	STOP_BM := directionToBitMap(DirStop)
	UP_BM := directionToBitMap(DirUp)
	DOWN_BM := directionToBitMap(DirDown)

	if drivingDirection == DirUp {
		// Search for orders above the elevator upwards first, then downwards, then up again
		for floor := lastFloor + 1; floor < N_FLOORS; floor++ {
			if orderHandler.floorState[floor] & (STOP_BM | UP_BM) > 0 {
				return floor
			}
		}

		for floor := N_FLOORS - 1; floor >= 0; floor-- {
			if orderHandler.floorState[floor] & (STOP_BM | DOWN_BM) > 0 {
				return floor
			}
		}

		for floor := 0; floor < lastFloor; floor++ {
			if orderHandler.floorState[floor] & (STOP_BM | UP_BM) > 0 {
				return floor
			}
		}

	} else {
		// Same as above, but backwards
		for floor := lastFloor - 1; floor >= 0; floor-- {
			if orderHandler.floorState[floor] & (STOP_BM | DOWN_BM) > 0 {
				return floor
			}
		}

		for floor := 0; floor < N_FLOORS; floor++ {
			if orderHandler.floorState[floor] & (STOP_BM | UP_BM) > 0 {
				return floor
			}
		}

		for floor := N_FLOORS - 1; floor >= lastFloor; floor-- {
			if orderHandler.floorState[floor] & (STOP_BM | DOWN_BM) > 0 {
				return floor
			}
		}
	}
	return -1
}

func stoppedAtFloor(orderHandler *OrderHandler, floor int, nextOrderFloor int) { // Theodor fix, lampene for opp og ned skal ikke slukkes hver gang heisen stopper
	orderHandler.floorState[floor] = 0
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
	if nextOrderFloor > floor {
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
	} else if nextOrderFloor < floor {								// BUG: lamp isnt turned off, as order is considered fulfilled when it shouldnt be
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
	}
}

func clearAllOrders(orderHandler *OrderHandler) {
	for floor := range N_FLOORS {
		elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
		orderHandler.floorState[floor] = 0
	}
}