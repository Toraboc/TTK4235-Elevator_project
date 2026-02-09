package main

import "time"

const numberOfFloors = 4

type NodeId [4]byte // Change to int32

const (
	// port is the UDP port used for both listening and broadcasting.
	port = 42067
	// broadcast is the IPv4 broadcast address and port used for discovery.
	broadcastAddress = "255.255.255.255"
	// sendHz is the broadcast frequency in Hz.
	sendHz = 5
	// printEvery is the interval used for peer list logging.
	printHz = 1
	// staleThreshold is the duration after which a peer is considered stale.
	staleThreshold = 500 * time.Millisecond
)
