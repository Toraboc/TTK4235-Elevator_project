package main

import (
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"sort"
	"time"
)

var nodesOnline NetworkState
var broadcast string

/*
TODO LIST:
- Fix knowsMe structure to work properly
- implement nodesOnline structure to keep track of which nodes are online
*/

// networkProcess starts the UDP listener and broadcaster for network communication.
func networkProcess() {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %d\n", getOwnId())
	nodesOnline = NetworkState{}
	knowsMe.node = make(map[NodeId]bool)

	go udpListen()
	udpBroadcast()

}

// getOwnId returns the IPv4 address of the computer as a NodeId.
func getOwnId() NodeId {
	var id NodeId
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return id
	}
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP == nil {
			continue
		}
		ip := ipnet.IP.To4()
		if ip == nil || ip.IsLoopback() {
			continue
		}
		copy(id[:], ip)
		return id
	}
	return id
}

// createOutgoingSync constructs a SyncMessage representing the current worldview.
func createOutgoingSync() SyncMessage {
	worldviewMutex.Lock()
	defer worldviewMutex.Unlock()

	syncMsg := SyncMessage{}
	syncMsg.Id = getOwnId()
	syncMsg.Orders = worldview.orders
	syncMsg.MyState = worldview.elevatorStates[syncMsg.Id]
	syncMsg.KnownNodes = make([]NodeId, len(nodesOnline.connectedNodes))
	for i, node := range nodesOnline.connectedNodes {
		syncMsg.KnownNodes[i] = node.id
	}
	return syncMsg
}

// udpBroadcast continuously broadcasts the SyncMessage over UDP at the configured sendHz.
func udpBroadcast() {
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(broadcastAddress), Port: port})
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	sendTimer := time.NewTicker(time.Second / sendHz)
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

// newKnownNodeSet creates an initialized nodeSet.
func newKnownNodeSet() *KnownNodeSet {
	return &KnownNodeSet{lastSeen: make(map[string]time.Time)}
}

// seen records that the given IP was observed now.
func (nodeSet *KnownNodeSet) nodeSeen(ip string) {
	nodeSet.mu.Lock()
	nodeSet.lastSeen[ip] = time.Now()
	nodeSet.mu.Unlock()
}

// list returns the sorted list of active peer IPs and prunes stale entries.
func (nodes *KnownNodeSet) listActivePeers() []string {
	nodes.mu.Lock()
	defer nodes.mu.Unlock()
	now := time.Now()
	for ip, t := range nodes.lastSeen {
		if now.Sub(t) > staleThresholdMs*time.Millisecond {
			delete(nodes.lastSeen, ip)
		}
	}
	ips := make([]string, 0, len(nodes.lastSeen))
	for ip := range nodes.lastSeen {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	//fmt.Println("Active peers:", ips)
	return ips
}

// udpListen listens for incoming SyncMessages over UDP, updates the nodeSet, and calls mergeWorldview on each received message.
func udpListen() {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return
	}
	defer conn.Close()

	peers := newKnownNodeSet()
	go func() {
		printTicker := time.NewTicker(time.Second / 100) // last number controls how often inactive peers are pruned (Hz)
		defer printTicker.Stop()
		for range printTicker.C {
			peers.listActivePeers()
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

		ip := fmt.Sprintf("%d.%d.%d.%d", syncMsg.Id[0], syncMsg.Id[1], syncMsg.Id[2], syncMsg.Id[3])
		peers.nodeSeen(ip)

		syncMsg.updateKnowsMe()
		//fmt.Printf("Knows about me %v\n", knowsMe.node)
		mergeWorldView(syncMsg)
	}
}

// updateKnowsMe updates the knowsMe structure based on the received SyncMessage.
func (syncMsg SyncMessage) updateKnowsMe() { // This is chatted, ignore
	knowsMe.mu.Lock()
	defer knowsMe.mu.Unlock()

	myID := getOwnId()
	knowsMe2 := slices.Contains(syncMsg.KnownNodes, myID)
	knowsMe.node[syncMsg.Id] = knowsMe2
}
