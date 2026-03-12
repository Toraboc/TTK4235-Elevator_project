package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	. "project/orderHandler"
	. "project/shared"
)

func NetworkProcess(orderHandler *OrderHandler, ConnectedNodesUpdateChannel chan NodeIdSet, WorldViewMergeChannel chan SyncView) {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %v\n", GetMyId())
	knownNodes := newKnownNodes()
	nodesAwareOfMe := newNodesAwareOfMe()

	go printConnectedNodes(knownNodes, nodesAwareOfMe) // For Debugging
	go pruneNodes(knownNodes, nodesAwareOfMe, ConnectedNodesUpdateChannel)
	go udpListen(knownNodes, nodesAwareOfMe, WorldViewMergeChannel)
	udpBroadcast(orderHandler, knownNodes)
}

func udpBroadcast(orderHandler *OrderHandler, KnownNodes *KnownNodes) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(BroadcastAddress), Port: Port})
	if err != nil {
		panic("Failed to dial UDP: " + err.Error())
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

func udpListen(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe, WorldViewMergeChannel chan SyncView) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: Port})
	if err != nil {
		panic("Failed to listen on UDP: " + err.Error())
	}
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			panic("Failed to read from UDP: " + err.Error())
		}
		var syncMsg SyncMessage
		err = json.Unmarshal(buf[:n], &syncMsg)
		if err != nil {
			continue
		}

		knownNodes.nodeSeen(syncMsg.Id)
		nodesAwareOfMe.update(syncMsg)
		WorldViewMergeChannel <- SyncView{NodeId: syncMsg.Id, ElevatorState: syncMsg.MyState, Orders: syncMsg.Orders}
	}
}
