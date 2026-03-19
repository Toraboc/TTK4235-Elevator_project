package orderHandler

import (
	"fmt"
	. "project/shared"
)

func OrderProcess(channels OrderHandlerInterface) {
	fmt.Println("Starting order process")
	worldView := newWorldView()

	for {
		select {
		case connectedNodes := <-channels.ConnectedNodesUpdateCh:
			worldView.ConnectedNodes = connectedNodes.Clone()
			updateTargetFloorIfChanged(channels, &worldView)

		case syncView := <-channels.SyncMergeCh:
			worldView.merge(syncView.NodeId, syncView.ElevatorState, syncView.Orders)
			updateTargetFloorIfChanged(channels, &worldView)

		case elevatorState := <-channels.ElevatorStateCh:
			worldView.ElevatorStates[GetMyId()] = elevatorState
			updateTargetFloorIfChanged(channels, &worldView)

		case orderCompleted := <-channels.OrderCompletedCh:
			worldView.completedOrder(orderCompleted)
			newTargetFloor, changed, err := updateTargetFloorIfChanged(channels, &worldView)
			if err != nil {
				panic(err.Error())
			}

			// This means that we still have an order in the opposite direction at this floor.
			// Then we also need to take this order afterwards.
			if !changed && newTargetFloor == orderCompleted.Floor {
				channels.TargetFloorCh <- newTargetFloor
			}

		case newOrder := <-channels.NewOrderCh:
			worldView.newOrder(newOrder)
			updateTargetFloorIfChanged(channels, &worldView)

		case responseCh := <-channels.RequestSyncCh:
			var syncData SyncData
			syncData.NodeId = GetMyId()
			syncData.Orders = *worldView.Orders[syncData.NodeId].Clone()
			syncData.ElevatorState = worldView.ElevatorStates[syncData.NodeId]
			responseCh <- syncData
		}
	}
}

func updateTargetFloorIfChanged(channels OrderHandlerInterface, worldView *WorldView) (int, bool, error) {
	targetFloor, changed, err := worldView.handleStateChange()
	myId := GetMyId()
	channels.ConfirmedOrdersCh <- getConfirmedOrders(worldView.Orders[myId], myId)
	if err == nil && changed {
		channels.TargetFloorCh <- targetFloor
	}
	return targetFloor, changed, err
}

