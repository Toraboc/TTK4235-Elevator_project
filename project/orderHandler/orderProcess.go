package orderHandler

import (
	"fmt"
	. "project/shared"
)

func pushTargetFloorIfChanged(channels OrderChannels, worldView *WorldView) {
	targetFloor, changed, err := worldView.handleStateChange()
	if err == nil && changed {
		channels.TargetFloorCh <- targetFloor
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
			myId := GetMyId()
			myOrders := worldView.Orders[myId]
			myOrders.CabOrders[myId][orderCompleted.Floor] = FINISHED
			switch orderCompleted.Direction {
			case UP:
				myOrders.HallUpOrders[orderCompleted.Floor] = FINISHED
			case DOWN:
				myOrders.HallDownOrders[orderCompleted.Floor] = FINISHED
			}
			worldView.Orders[myId] = myOrders
			pushTargetFloorIfChanged(channels, &worldView)

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

		case responseCh := <-channels.ConfirmedOrdersReqCh:
			responseCh <- worldView.getConfirmedOrders()
		}
	}
}
