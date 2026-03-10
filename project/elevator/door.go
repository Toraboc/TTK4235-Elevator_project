package elevator

import (
	"errors"
	"github.com/angrycompany16/driver-go/elevio"
	"time"
)

const DoorOpenTime = 3 * time.Second

type Door struct {
	isOpen bool
	timer  *time.Timer
}

func (door *Door) IsOpen() bool {
	return door.isOpen
}

func (door *Door) Open() {
	elevio.SetDoorOpenLamp(true)
	door.isOpen = true

	if door.timer == nil {
		door.timer = time.NewTimer(DoorOpenTime)
	}

	door.timer.Reset(DoorOpenTime)
}

func (door *Door) IsObstructed() bool {
	return elevio.GetObstruction()
}

func (door *Door) Close() error {
	if door.IsObstructed() {
		return errors.New("The door is obstructed, cannot close the door")
	}

	elevio.SetDoorOpenLamp(false)
	door.isOpen = false

	if door.timer == nil {
		door.timer = time.NewTimer(DoorOpenTime)
	}

	door.timer.Stop()

	return nil
}

func (door *Door) CloseTrigger() <-chan time.Time {
	if door.timer == nil {
		door.timer = time.NewTimer(DoorOpenTime)
		door.timer.Stop()
	}

	return door.timer.C
}
