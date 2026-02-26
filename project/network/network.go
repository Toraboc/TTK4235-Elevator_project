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
	- Add sleep to udpListen??
*/

// NetworkProcess starts the UDP listener and broadcaster for network communication.
func NetworkProcess() {
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

	go pruneTicker(knownNodes, nodesAwareOfMe)
	go udpListen(knownNodes, nodesAwareOfMe)
	udpBroadcast(knownNodes)
}

// createOutgoingSync constructs a SyncMessage representing the current worldview.
func createOutgoingSync(knownNodes *KnownNodes) SyncMessage {
	worldview := GetWorldView()

	syncMsg := SyncMessage{}
	syncMsg.Id = GetMyId()
	syncMsg.Orders = worldview.Orders
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
func udpBroadcast(KnownNodes *KnownNodes) {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(BroadcastAddress), Port: Port})
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	sendTimer := time.NewTicker(time.Second / SendHz)
	defer sendTimer.Stop()

	for range sendTimer.C {
		syncMsg := createOutgoingSync(KnownNodes)
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
func udpListen(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
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
		MergeWorldView(syncMsg)
	}
}

// pruneTicker periodically prunes stale nodes from knownNodes and nodesAwareOfMe, and updates the connected nodes.
func pruneTicker(knownNodes *KnownNodes, nodesAwareOfMe *NodesAwareOfMe) {
	ticker := time.NewTicker(time.Second / 100) // last number controls how often inactive peers are pruned (Hz)
	defer ticker.Stop()
	for range ticker.C {
		knownNodes.pruneStale()
		nodesAwareOfMe.pruneStale()
		UpdateConnectedNodes(GetConnectedNodes(knownNodes, nodesAwareOfMe))
	}
}
