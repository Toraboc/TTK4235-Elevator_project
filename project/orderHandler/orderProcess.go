package orderHandler

import (
	"fmt"
	. "project/shared"
)

func pushTargetFloorIfChanged(channels OrderChannels, worldView *WorldView) (int, bool, error) {
	targetFloor, changed, err := worldView.handleStateChange()
	channels.ConfirmedOrdersCh <- worldView.getConfirmedOrders()
	if err == nil && changed {
		channels.TargetFloorCh <- targetFloor
	}
	return targetFloor, changed, err
}

func handleOrderCompleted(channels OrderChannels, worldView *WorldView, orderCompleted OrderCompleted) {
	myId := GetMyId()
	myOrders := worldView.Orders[myId]

	myOrders.CabOrders[myId][orderCompleted.Floor] = FINISHED

	var hadHallOrder bool
	switch orderCompleted.Direction {
	case UP:
		hadHallOrder = myOrders.HallUpOrders[orderCompleted.Floor] == CONFIRMED
		if hadHallOrder {
			myOrders.HallUpOrders[orderCompleted.Floor] = FINISHED
		}
	case DOWN:
		hadHallOrder = myOrders.HallDownOrders[orderCompleted.Floor] == CONFIRMED
		if hadHallOrder {
			myOrders.HallDownOrders[orderCompleted.Floor] = FINISHED
		}
	}

	if !hadHallOrder {
		worldView.hallRequestAssigner()
		targetFloor, err := worldView.getNextTargetFloor()
		if err != nil {
			panic(err.Error())
		}

		if targetFloor == orderCompleted.Floor {
			switch orderCompleted.Direction {
			case UP:
				myOrders.HallDownOrders[orderCompleted.Floor] = FINISHED
			case DOWN:
				myOrders.HallUpOrders[orderCompleted.Floor] = FINISHED
			}

		}
	}

	newTargetFloor, changed, err := pushTargetFloorIfChanged(channels, worldView)
	if err != nil {
		panic(err.Error())
	}

	// This means that we still have an order in the opposite direction at this floor.
	// Then we also need to take this order afterwards.
	if !changed && newTargetFloor == orderCompleted.Floor {
		channels.TargetFloorCh <- newTargetFloor
	}
}

func OrderProcess(channels OrderChannels) {
	fmt.Println("Starting order process")
	worldView := newWorldView()

	for {
		select {
		case connectedNodes := <-channels.ConnectedNodesUpdateCh:
			worldView.ConnectedNodes = connectedNodes.Clone()
			pushTargetFloorIfChanged(channels, &worldView)

		case syncView := <-channels.WorldViewMergeCh:
			worldView.merge(syncView.NodeId, syncView.ElevatorState, syncView.Orders)
			pushTargetFloorIfChanged(channels, &worldView)

		case elevatorState := <-channels.ElevatorStateCh:
			worldView.ElevatorStates[GetMyId()] = elevatorState
			pushTargetFloorIfChanged(channels, &worldView)

		case orderCompleted := <-channels.OrderCompletedCh:
			handleOrderCompleted(channels, &worldView, orderCompleted)
		case newOrder := <-channels.NewOrderCh:
			myId := GetMyId()
			myOrders := worldView.Orders[myId]
			switch newOrder.OrderType {
			case HALLUP:
				if myOrders.HallUpOrders[newOrder.Floor] == NO_ORDER {
					myOrders.HallUpOrders[newOrder.Floor] = UNCONFIRMED
				}
			case HALLDOWN:
				if myOrders.HallDownOrders[newOrder.Floor] == NO_ORDER {
					myOrders.HallDownOrders[newOrder.Floor] = UNCONFIRMED
				}
			case CAB:
				myCabOrders := myOrders.CabOrders[myId]
				if myCabOrders[newOrder.Floor] == NO_ORDER {
					myCabOrders[newOrder.Floor] = UNCONFIRMED
				}
				myOrders.CabOrders[myId] = myCabOrders
			}
			worldView.Orders[myId] = myOrders
			pushTargetFloorIfChanged(channels, &worldView)

		case responseCh := <-channels.WorldViewReqCh:
			responseCh <- worldView.clone()
		}
	}
}
