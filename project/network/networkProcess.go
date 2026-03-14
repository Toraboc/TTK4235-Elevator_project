package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	. "project/orderHandler"
	. "project/shared"
)

func NetworkProcess(channels OrderChannels) {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %v\n", GetMyId())
	knownNodes := newKnownNodes()
	nodesAwareOfMe := newNodesAwareOfMe()

	nodeUpdate(knownNodes, nodesAwareOfMe, channels.ConnectedNodesUpdateCh)
	go pruneNodes(knownNodes, nodesAwareOfMe, channels.ConnectedNodesUpdateCh)
	go udpListen(knownNodes, nodesAwareOfMe, channels.ConnectedNodesUpdateCh, channels.WorldViewMergeCh)
	udpBroadcast(knownNodes, channels.WorldViewReqCh)

}

func udpBroadcast(knownNodes *KnownNodes, worldViewReqCh WorldViewRequestCh) {
	var conn *net.UDPConn
	for {
		var err error
		conn, err = net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(BroadcastAddress), Port: Port})
		if err == nil {
			break
		}
		fmt.Println("Failed to dial UDP, retrying in 1 second:", err)
		time.Sleep(1 * time.Second)
	}
	defer conn.Close()

	sendTimer := time.NewTicker(time.Second / SendHz)
	defer sendTimer.Stop()

	for range sendTimer.C {
		syncMsg := createOutgoingSync(worldViewReqCh, knownNodes)
		data, err := json.Marshal(syncMsg)
		if err != nil {
			fmt.Println("Error marshaling sync message:", err)
			continue
		}
		_, err = conn.Write(data)
		if err != nil {
			// fmt.Println("Error writing to UDP:", err)
			continue
		}
	}
}

func udpListen(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe, connectedNodesUpdateCh chan<- NodeIdSet, worldViewMergeCh chan<- SyncView) {
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

		if syncMsg.Id == GetMyId() {
			continue
		}

		knownNodes.nodeSeen(syncMsg.Id, nodesAwareOfMe, connectedNodesUpdateCh)
		nodesAwareOfMe.update(syncMsg, knownNodes, connectedNodesUpdateCh)

		worldViewMergeCh <- SyncView{NodeId: syncMsg.Id, ElevatorState: syncMsg.MyState, Orders: syncMsg.Orders}
	}
}
