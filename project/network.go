package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"time"
)

func networkProcess() {
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %d\n", getOwnId())

	go udpListen()
	udpBroadcast()

}

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

	netState := NetworkState{}
	connectedNodes := netState.connectedNodes

	syncMsg := SyncMessage{}
	syncMsg.id = getOwnId()
	syncMsg.orders = worldview.orders
	syncMsg.myState = worldview.elevatorStates[syncMsg.id]
	syncMsg.knownNodes = make([]NodeId, len(connectedNodes))
	for i, node := range connectedNodes {
		syncMsg.knownNodes[i] = node.id
	}
	return syncMsg
}

// udpBroadcast continuously broadcasts the SyncMessage over UDP at the configured sendHz.
func udpBroadcast() {
	var broadcast = fmt.Sprintf("%s:%d", broadcastAddress, port)
	conn, err := net.DialUDP("udp4", nil, &net.UDPAddr{IP: net.ParseIP(broadcastAddress), Port: port})
	if err != nil {
		return
	}
	defer conn.Close()

	sendTimer := time.NewTicker(time.Second / sendHz)
	defer sendTimer.Stop()

	broadcastAddr, _ := net.ResolveUDPAddr("udp4", broadcast)

	for range sendTimer.C {
		syncMsg := createOutgoingSync()
		data, err := json.Marshal(syncMsg)
		if err != nil {
			return
		}
		_, _ = conn.WriteToUDP(data, broadcastAddr)
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
	// drop stale nodes
	now := time.Now()
	for ip, t := range nodes.lastSeen {
		if now.Sub(t) > staleThreshold*time.Millisecond {
			delete(nodes.lastSeen, ip)
		}
	}
	ips := make([]string, 0, len(nodes.lastSeen))
	for ip := range nodes.lastSeen {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	return ips
}

// udpListen listens for incoming SyncMessages over UDP, updates the nodeSet, and calls mergeWorldview on each received message.
func udpListen() {
	// Listen for incoming UDP packets on the specified port.
	// update nodeset when receiving packets
	// call mergeWorldview on received packets
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return
	}
	defer conn.Close()

	peers := newKnownNodeSet()
	go peers.listActivePeers()

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

		fmt.Printf("Received sync from %d\n", syncMsg.id)
		peers.nodeSeen(fmt.Sprintf("%d.%d.%d.%d", syncMsg.id[0], syncMsg.id[1], syncMsg.id[2], syncMsg.id[3]))
		mergeWorldView(syncMsg)

	}

}
