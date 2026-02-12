package elevator

import (
	"Driver-go/elevio"
	"project/network"
	"project/orders"
	. "project/shared"
	"time"
)

type LightStrip [NumberOfFloors]bool

type LightStatus struct {
	hallUp LightStrip
	hallDown LightStrip
	cab LightStrip
}

var lightStatus LightStatus

func (lightStrip *LightStrip) update(lamp elevio.ButtonType, confirmedOrders [NumberOfFloors]OrderStatus) {
	for floor := range NumberOfFloors {
		orderConfirmed := confirmedOrders[floor] == NEW
		if orderConfirmed != lightStrip[floor] {
			lightStrip[floor] = orderConfirmed
			elevio.SetButtonLamp(lamp, floor, orderConfirmed)
		}
	}
}

func (lightStatus *LightStatus) updateLights(confirmedOrders orders.ConfirmedOrders) {
	ownId := network.GetOwnId() // TODO: Use a cache for this
	cabOrders := confirmedOrders.Cab[ownId]

	lightStatus.hallUp.update(elevio.BT_HallUp, confirmedOrders.HallUp)
	lightStatus.hallDown.update(elevio.BT_HallDown, confirmedOrders.HallDown)
	lightStatus.cab.update(elevio.BT_Cab, cabOrders)
}

func (lightStatus *LightStatus) Init() {
	for floor := range NumberOfFloors {
		elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
		elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
		elevio.SetButtonLamp(elevio.BT_Cab, floor, false)

		lightStatus.hallUp[floor] = false
		lightStatus.hallDown[floor] = false
		lightStatus.cab[floor] = false
	}
}

func handleLights() {
	lightStatus.Init()
	for {
		time.Sleep(40 * time.Millisecond)

		ownId := network.GetOwnId() // TODO: This will be removed then this stops something we need to pass to the function below
		confirmedOrders := orders.GetConfirmedOrders(ownId)
		lightStatus.updateLights(confirmedOrders)
	}
}
