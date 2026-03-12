package shared

import (
	"fmt"
	"net"
	"strings"
)

type NodeId uint32

type NodeIdSet map[NodeId]struct{}

var myId NodeId

func (set NodeIdSet) Contains(nodeId NodeId) bool {
	_, exists := set[nodeId]
	return exists
}

func (set NodeIdSet) Add(nodeId NodeId) {
	set[nodeId] = struct{}{}
}

func (set NodeIdSet) Concat(other NodeIdSet) {
	for nodeId := range other {
		set.Add(nodeId)
	}
}

func (set NodeIdSet) Clone() NodeIdSet {
	newSet := make(NodeIdSet)

	for key, _ := range set {
		newSet[key] = struct{}{}
	}

	return newSet
}

func (set NodeIdSet) Remove(nodeId NodeId) {
	delete(set, nodeId)
}

func (set NodeIdSet) String() string {
	var builder strings.Builder

	builder.WriteString("[")

	for nodeId := range set {
		builder.WriteString(nodeId.String())
		builder.WriteString(", ")
	}

	builder.WriteString("]")
	return builder.String()
}

func NewNodeIdSet(nodeIds []NodeId) NodeIdSet {
	set := make(NodeIdSet)
	for _, node := range nodeIds {
		set[node] = struct{}{}
	}
	return set
}

func (id NodeId) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(id>>24), byte(id>>16), byte(id>>8), byte(id))
}

func NodeIdListToStrings(ids []NodeId) []string {
	result := make([]string, len(ids))
	for i, id := range ids {
		result[i] = id.String()
	}
	return result
}

// getIpAddress returns the IPv4 address of the computer as a NodeId. Heavy process, should only be called once at startup. If no valid IP is found, returns 0.
func getIpAddress() NodeId {
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
	if myId == 0 {
		myId = getIpAddress()
	}
	return myId
}
