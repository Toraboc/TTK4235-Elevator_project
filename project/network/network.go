package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	. "project/orders"
	. "project/shared"
)

/*
	TODO:
	- Brew coffee
	- Take a nap
	- You know, the usual
*/

// NetworkProcess starts the UDP listener and broadcaster for network communication.
func NetworkProcess() {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %s\n", NodeIdtoString(GetMyId()))
	nodesAwareOfMe := newNodesAwareOfMe()
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Knows about me:", nodesAwareOfMe.KnowsMe()) // DEBUG
		}
	}()

	go udpListen(nodesAwareOfMe)
	udpBroadcast()
}

func NodeIdtoString(id NodeId) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(id>>24), byte(id>>16), byte(id>>8), byte(id))
}

func NodeIdListToStrings(ids []NodeId) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = NodeIdtoString(id)
	}
	return result
}

// createOutgoingSync constructs a SyncMessage representing the current worldview.
func createOutgoingSync() SyncMessage {
	worldview := GetWorldView()

	syncMsg := SyncMessage{}
	syncMsg.Id = GetMyId()
	syncMsg.Orders = worldview.Orders
	syncMsg.MyState = worldview.ElevatorStates[syncMsg.Id]
	syncMsg.KnownNodes = make([]NodeId, len(worldview.ConnectedNodes))
	copy(syncMsg.KnownNodes, worldview.ConnectedNodes)
	syncMsg.SendTime = time.Now()
	return syncMsg
}

// udpBroadcast continuously broadcasts the SyncMessage over UDP at the configured sendHz.
func udpBroadcast() {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(BroadcastAddress), Port: Port})
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	sendTimer := time.NewTicker(time.Second / SendHz)
	defer sendTimer.Stop()

	for range sendTimer.C {
		syncMsg := createOutgoingSync()
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
func udpListen(nodesAwareOfMe *NodesAwareOfMe) {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: Port})
	if err != nil {
		return
	}
	defer conn.Close()

	peers := newKnownNodes()
	go func() {
		printTicker := time.NewTicker(time.Second / 100) // last number controls how often inactive peers are pruned (Hz)
		defer printTicker.Stop()
		for range printTicker.C {
			peers.updateConnectedNodes()
			nodesAwareOfMe.purgeStale()
		}
	}()

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

		ip := syncMsg.Id
		peers.nodeSeen(ip)

		nodesAwareOfMe.update(syncMsg)
		//nodesAwareOfMe.Print() // DEBUG
		MergeWorldView(syncMsg)
	}
}
