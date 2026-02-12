package network

import (
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"sort"
	"sync"
	"time"

	. "project/orders"
	. "project/shared"
)

type NetworkState struct {
	ConnectedNodes []NetworkNode
}

type KnowsMe struct {
	Node map[NodeId]bool
	Mu   sync.Mutex
}

type KnownNodes struct {
	Mu       sync.Mutex
	LastSeen map[string]time.Time
}

var nodesOnline NetworkState
var broadcast string
var knowsMe KnowsMe

/*
TODO LIST:
- Fix knowsMe structure to work properly
- implement nodesOnline structure to keep track of which nodes are online
- change to 32bit NodeId
- Implement clock offset compensation
*/

// NetworkProcess starts the UDP listener and broadcaster for network communication.
func NetworkProcess() {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %d\n", getOwnId())
	nodesOnline = NetworkState{}
	knowsMe.Node = make(map[NodeId]bool)

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

// createOutgoingSync constructs a SyncMessage representing the current worldView.
func createOutgoingSync() SyncMessage {
	worldView := GetWorldView()

	syncMsg := SyncMessage{}
	syncMsg.Id = getOwnId()
	syncMsg.Orders = worldView.Orders
	syncMsg.MyState = worldView.ElevatorStates[syncMsg.Id]
	syncMsg.KnownNodes = make([]NodeId, len(nodesOnline.ConnectedNodes))
	for i, node := range nodesOnline.ConnectedNodes {
		syncMsg.KnownNodes[i] = node.Id
	}
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

// newKnownNodes creates an initialized KnownNodes.
func newKnownNodes() *KnownNodes {
	return &KnownNodes{LastSeen: make(map[string]time.Time)}
}

// seen records that the given IP was observed now.
func (nodeSet *KnownNodes) nodeSeen(ip string) {
	nodeSet.Mu.Lock()
	nodeSet.LastSeen[ip] = time.Now()
	nodeSet.Mu.Unlock()
}

// list returns the sorted list of active peer IPs and prunes stale entries.
func (nodes *KnownNodes) listActivePeers() []string {
	nodes.Mu.Lock()
	defer nodes.Mu.Unlock()
	for ip, seenAt := range nodes.LastSeen {
		if time.Since(seenAt) > StaleThreshold {
			delete(nodes.LastSeen, ip)
		}
	}
	ips := make([]string, 0, len(nodes.LastSeen))
	for ip := range nodes.LastSeen {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	//fmt.Println("Active peers:", ips)
	return ips
}

// udpListen listens for incoming SyncMessages over UDP, updates the nodeSet, and calls mergeWorldView on each received message.
func udpListen() {
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

		updateKnowsMe(syncMsg)
		//fmt.Printf("Knows about me %v\n", knowsMe.node)
		MergeWorldView(syncMsg)
	}
}

// updateKnowsMe updates the knowsMe structure based on the received SyncMessage.
func updateKnowsMe(syncMsg SyncMessage) { // This is chatted, ignore
	knowsMe.Mu.Lock()
	defer knowsMe.Mu.Unlock()

	myID := getOwnId()
	knowsMe2 := slices.Contains(syncMsg.KnownNodes, myID)
	knowsMe.Node[syncMsg.Id] = knowsMe2
}
