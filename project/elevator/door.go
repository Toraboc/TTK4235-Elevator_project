package elevator

import (
	"Driver-go/elevio"
	"errors"
	"time"
)

const doorOpenTime = 3 * time.Second

type Door struct {
	isOpen     bool
	changeTime time.Time
}

func (door Door) IsOpen() bool {
	return door.isOpen
}

func (door Door) Open() {
	elevio.SetDoorOpenLamp(true)
	door.isOpen = true
	door.changeTime = time.Now()
}

func (door Door) IsObsrcted() bool {
	return elevio.GetObstruction()
}

func (door Door) Close() error {
	if door.IsObsrcted() {
		return errors.New("The door is obstrcted, cannot close the door")
	}

	elevio.SetDoorOpenLamp(false)
	door.isOpen = false
	door.changeTime = time.Now()

	return nil
}

func (door Door) ShouldClose() bool {
	return !door.IsOpen() && time.Since(door.changeTime) > doorOpenTime
}
