package network

import "time"

const (
	// UDP port used for both listening and broadcasting.
	Port = 44043
	// IPv4 broadcast address and port used for discovery.
	BroadcastAddress = "255.255.255.255"
	// broadcast frequency in Hz.
	SendFrequency = 250
	// frequency for pruning stale peers.
	updateConnectedNodesFrequency = 10
	// duration after which a peer is considered stale.
	StaleThreshold = 200 * time.Millisecond
)
