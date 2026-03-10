package elevator

import (
	"time"
)

const ElevatorServer = "localhost:15657"
const TimeBetweenFloors = 4 * time.Second
const PositionPollInterval = 50 * time.Millisecond
const DoorOpenTime = 3 * time.Second
