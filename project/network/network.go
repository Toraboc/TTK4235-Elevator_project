package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	. "project/orderHandler"
	. "project/shared"
)

/*
	TODO:
	- Brew coffee
	- Take a nap
	- You know, the usual
	- Add sleep to udpListen??
*/

// NetworkProcess starts the UDP listener and broadcaster for network communication.
func NetworkProcess(orderHandler *OrderHandler) {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %v\n", GetMyId())
	knownNodes := newKnownNodes()
	nodesAwareOfMe := newNodesAwareOfMe()
	go func() { // Debug loop to print known nodes and nodes aware of me every second
		for {
			time.Sleep(1 * time.Second)
			knownNodes.Print()
			nodesAwareOfMe.Print()
		}
	}()

	go udpListen(orderHandler, knownNodes, nodesAwareOfMe)
	udpBroadcast(orderHandler, knownNodes)
	go pruneNodes(orderHandler,knownNodes, nodesAwareOfMe)
	go udpListen(orderHandler, knownNodes, nodesAwareOfMe)
	udpBroadcast(orderHandler, knownNodes)
}

// createOutgoingSync constructs a SyncMessage representing the current worldview.
func createOutgoingSync(orderHandler *OrderHandler, knownNodes *KnownNodes) SyncMessage {
	worldview := orderHandler.GetWorldView()

	syncMsg := SyncMessage{}
	syncMsg.Id = GetMyId()
	syncMsg.Orders  = worldview.Orders[syncMsg.Id]
	syncMsg.MyState = worldview.ElevatorStates[syncMsg.Id]
	knownNodes.mu.Lock()
	defer knownNodes.mu.Unlock()
	syncMsg.KnownNodes = make([]NodeId, 0, len(knownNodes.LastSeen))
	for id := range knownNodes.LastSeen {
		syncMsg.KnownNodes = append(syncMsg.KnownNodes, id)
	}
	syncMsg.SendTime = time.Now()
	return syncMsg
}

// udpBroadcast continuously broadcasts the SyncMessage over UDP at the configured sendHz.
func udpBroadcast(orderHandler *OrderHandler, KnownNodes *KnownNodes) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(BroadcastAddress), Port: Port})
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	sendTimer := time.NewTicker(time.Second / SendHz)
	defer sendTimer.Stop()

	for range sendTimer.C {
		syncMsg := createOutgoingSync(orderHandler, KnownNodes)
		data, err := json.Marshal(syncMsg)
		if err != nil {
			fmt.Println("Error marshaling sync message:", err)
			continue
		}
		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error writing to UDP:", err)
			continue
		}
	}
}

// udpListen listens for incoming SyncMessages over UDP, updates the nodeSet, and calls mergeWorldview on each received message.
func udpListen(orderHandler *OrderHandler, knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: Port})
	if err != nil {
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		var syncMsg SyncMessage
		err = json.Unmarshal(buf[:n], &syncMsg)
		if err != nil {
			continue
		}

		knownNodes.nodeSeen(syncMsg.Id)
		nodesAwareOfMe.update(syncMsg)
		orderHandler.MergeWorldView(syncMsg.Id, syncMsg.MyState, syncMsg.Orders)
	}
}
