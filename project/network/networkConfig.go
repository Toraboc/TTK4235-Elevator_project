package network

import "time"

const (
	// port is the UDP port used for both listening and broadcasting.
	Port = 42067
	// broadcast is the IPv4 broadcast address and port used for discovery.
	BroadcastAddress = "255.255.255.255"
	// sendHz is the broadcast frequency in Hz.
	SendHz = 5
	// printEvery is the interval used for peer list logging.
	PrintHz = 1
	// staleThreshold is the duration after which a peer is considered stale.
	StaleThreshold = 500 * time.Millisecond
)
