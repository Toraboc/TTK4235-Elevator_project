package elevator

import (
	"time"
)

const elevatorServer = "localhost:15657"
const timeBetweenFloors = 4 * time.Second
const positionPollInterval = 50 * time.Millisecond
const doorOpenTime = 3 * time.Second
