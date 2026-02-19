package network

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	. "project/orders"
	. "project/shared"
)

var myId NodeId

/*
	TODO:
*/

// NetworkProcess starts the UDP listener and broadcaster for network communication.
func NetworkProcess() {
	myId = getOwnId()
	fmt.Println("Starting network process")
	fmt.Printf("My Ip: %s\n", NodeIdtoString(myId))
	otherNodes := PeersAwareOfMe{knowsAboutMe: make(map[NodeId]KnowsAboutMe)}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			fmt.Println("Knows about me:", GetKnowsAboutMe(&otherNodes)) // DEBUG
		}
	}()

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

// getOwnId returns the IPv4 address of the computer as a NodeId. Heavy process, should only be called once at startup. If no valid IP is found, returns 0.
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
		id32 := (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3])
		return NodeId(id32)
	}
	return 0
}

// GetMyId returns the NodeId of this node.
func GetMyId() NodeId {
	return myId
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

// Clock offset compensation adjusts the timestamps in the SyncMessage to account for clock differences between nodes.
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
		//printKnowsAboutMe(otherNodes) // DEBUG
		//clockOffsetCompensation(&syncMsg)
		MergeWorldView(syncMsg)
	}
}
