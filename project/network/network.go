package network

import (
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	. "project/orders"
	. "project/shared"
)

type KnowsAboutMe struct {
	Node         bool
	LastReceived time.Time
}

type PeersAwareOfMe struct {
	mu           sync.Mutex
	knowsAboutMe map[NodeId]KnowsAboutMe
}

type KnownNodes struct {
	Mu       sync.Mutex
	LastSeen map[NodeId]time.Time
}

var myId NodeId

/*
	TODO:
	- Implement a function that returns all the nodes that know about me as a nodeIdSet
*/

// NetworkProcess starts the UDP listener and broadcaster for network communication.
func NetworkProcess() {
	myId = GetOwnId()
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %s\n", NodeIdtoString(myId))
	otherNodes := PeersAwareOfMe{knowsAboutMe: make(map[NodeId]KnowsAboutMe)}

	go udpListen(&otherNodes)
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

// getOwnId returns the IPv4 address of the computer as a NodeId.
func GetOwnId() NodeId {
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
		id32 := (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
		return NodeId(id32)
	}
	return 0
}

// createOutgoingSync constructs a SyncMessage representing the current worldview.
func createOutgoingSync() SyncMessage {
	worldview := GetWorldView()

	syncMsg := SyncMessage{}
	syncMsg.Id = myId
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

// newKnownNodes creates an initialized KnownNodes.
func newKnownNodes() *KnownNodes {
	return &KnownNodes{LastSeen: make(map[NodeId]time.Time)}
}

// seen records that the given IP was observed now.
func (nodeSet *KnownNodes) nodeSeen(id NodeId) {
	nodeSet.Mu.Lock()
	nodeSet.LastSeen[id] = time.Now()
	nodeSet.Mu.Unlock()
}

// list returns the sorted list of active peer IPs and prunes stale entries.
func (nodes *KnownNodes) listActivePeers() []NodeId {
	nodes.Mu.Lock()
	defer nodes.Mu.Unlock()
	for id, seenAt := range nodes.LastSeen {
		if time.Since(seenAt) > StaleThreshold {
			delete(nodes.LastSeen, id)
		}
	}
	ids := make([]NodeId, 0, len(nodes.LastSeen))
	for id := range nodes.LastSeen {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	fmt.Println("Active peers:", NodeIdListToStrings(ids))

	UpdateConnectedNodes(ids)
	return ids
}

func clockOffsetCompensation(syncMsg *SyncMessage) {
	// This is a placeholder for clock offset compensation logic.
	offset := time.Since(syncMsg.SendTime)
	for order := range syncMsg.Orders.HallUpOrders {
		syncMsg.Orders.HallUpOrders[order].LastUpdate = syncMsg.Orders.HallUpOrders[order].LastUpdate.Add(offset)
	}
	for order := range syncMsg.Orders.HallDownOrders {
		syncMsg.Orders.HallDownOrders[order].LastUpdate = syncMsg.Orders.HallDownOrders[order].LastUpdate.Add(offset)
	}
	for nodeID, cabOrders := range syncMsg.Orders.CabOrders {
		for floor := range cabOrders {
			order := cabOrders[floor]
			order.LastUpdate = order.LastUpdate.Add(offset)
			cabOrders[floor] = order
		}
		syncMsg.Orders.CabOrders[nodeID] = cabOrders
	}
}

// udpListen listens for incoming SyncMessages over UDP, updates the nodeSet, and calls mergeWorldview on each received message.
func udpListen(otherNodes *PeersAwareOfMe) {
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
			purgeStaleKnowsMe(otherNodes)
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

		updateKnowsMe(syncMsg, otherNodes)
		fmt.Printf("Knows about me %v\n", otherNodes.knowsAboutMe)
		//clockOffsetCompensation(&syncMsg)
		MergeWorldView(syncMsg)
	}
}

// updateKnowsMe updates the knowsAboutMe based on the received SyncMessage.
func updateKnowsMe(syncMsg SyncMessage, otherNodes *PeersAwareOfMe) {
	otherNodes.mu.Lock()
	defer otherNodes.mu.Unlock()

	for i := range syncMsg.KnownNodes {
		if syncMsg.KnownNodes[i] == myId {
			entry := otherNodes.knowsAboutMe[syncMsg.Id]
			entry.Node = true
			entry.LastReceived = time.Now()
			otherNodes.knowsAboutMe[syncMsg.Id] = entry
		}
	}
}

func purgeStaleKnowsMe(peersAwareOfMe *PeersAwareOfMe) {
	peersAwareOfMe.mu.Lock()
	defer peersAwareOfMe.mu.Unlock()
	for id, entry := range peersAwareOfMe.knowsAboutMe {
		if time.Since(entry.LastReceived) > StaleThreshold {
			entry.Node = false
			peersAwareOfMe.knowsAboutMe[id] = entry
		}
	}
}
