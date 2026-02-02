package network

import (
	"fmt"
	"project/drivers/elevator"
	"project/drivers/orderHandler"
	"sync"
	"time"
	//"bytes"
	//"net"
	//"sort"
)

type NodeID int32

type NetworkNode struct {
	id      NodeID
	lastHB  time.Time
	knowsMe bool
	mu      sync.Mutex
}

type SyncMessage struct {
	orders     orderhandler.Orders
	myState    elevator.ElevatorState
	knownNodes []NodeID
}

type NetworkState struct {
	connectedNodes map[NodeID]NetworkNode
}

const (
	// port is the UDP port used for both listening and broadcasting.
	port = 42067
	// broadcast is the IPv4 broadcast address and port used for discovery.
	broadcast = "255.255.255.255:42067"
	// sendHz is the broadcast frequency in Hz.
	sendHz = 100
	// printEvery is the interval used for peer list logging.
	printEvery = time.Second
)

func udplisten() {

}

func udpbroadcast() {

}
