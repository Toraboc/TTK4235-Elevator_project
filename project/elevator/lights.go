package elevator

import (
	"github.com/angrycompany16/driver-go/elevio"
	//"project/network"
	"project/orders"
	. "project/shared"
	"time"
)

type LightType [NumberOfFloors]bool

type LightStatus struct {
	hallUp   LightType
	hallDown LightType
	cab      LightType
}

var lightStatus LightStatus

func (lightType *LightType) update(lamp elevio.ButtonType, confirmedOrders [NumberOfFloors]OrderStatus) {
	for floor := range NumberOfFloors {
		orderConfirmed := confirmedOrders[floor] == NEW
		if orderConfirmed != lightType[floor] {
			lightType[floor] = orderConfirmed
			elevio.SetButtonLamp(lamp, floor, orderConfirmed)
		}
	}
}

func (lightStatus *LightStatus) updateLights(confirmedOrders orders.ConfirmedOrders) {
	ownId := GetMyId() // TODO: Use a cache for this
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

		confirmedOrders := orders.GetConfirmedOrders()
		lightStatus.updateLights(confirmedOrders)
	}
}
