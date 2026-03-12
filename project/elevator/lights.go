package elevator

import (
	"github.com/angrycompany16/driver-go/elevio"
	//"project/network"
	. "project/orderHandler"
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

func (lightType *LightType) update(lamp elevio.ButtonType, confirmedOrders [NumberOfFloors]bool) {
	for floor := range NumberOfFloors {
		orderConfirmed := confirmedOrders[floor]
		if orderConfirmed != lightType[floor] {
			lightType[floor] = orderConfirmed
			elevio.SetButtonLamp(lamp, floor, orderConfirmed)
		}
	}
}

func (lightStatus *LightStatus) updateLights(confirmedOrders ConfirmedOrders) {
	lightStatus.hallUp.update(elevio.BT_HallUp, confirmedOrders.HallUp)
	lightStatus.hallDown.update(elevio.BT_HallDown, confirmedOrders.HallDown)
	lightStatus.cab.update(elevio.BT_Cab, confirmedOrders.Cab)
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

func handleLights(confirmedOrdersReqCh ConfirmedOrdersRequestCh) {
	lightStatus.Init()
	for {
		time.Sleep(40 * time.Millisecond)

		confirmedOrders := RequestConfirmedOrders(confirmedOrdersReqCh)
		lightStatus.updateLights(confirmedOrders)
	}
}
