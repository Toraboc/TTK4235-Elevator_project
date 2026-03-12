package network

import "time"

const (
	// IPv4 broadcast address and port used for discovery.
	BroadcastAddress = "255.255.255.255"
	// broadcast frequency in Hz.
	SendHz = 100
	// frequency for peer list logging.
	PrintHz = 2
	// frequency for pruning stale peers.
	PruneHz = 100
	// duration after which a peer is considered stale.
	StaleThreshold = 200 * time.Millisecond
)
