package elevator

import (
	"time"
)

const timeBetweenFloors = 4 * time.Second
const positionPollInterval = 50 * time.Millisecond
const doorOpenTime = 3 * time.Second
