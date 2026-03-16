package shared

import (
	"fmt"
	"maps"
	"net"
	"slices"
	"strings"
)

type NodeId uint32

type NodeIdSet map[NodeId]struct{}

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

	for key := range set {
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

	for i, nodeId := range slices.Sorted(maps.Keys(set)) {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(nodeId.String())
	}

	builder.WriteString("]")
	return builder.String()
}

func (id NodeId) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(id>>24), byte(id>>16), byte(id>>8), byte(id))
}

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
	panic("Failed to get IP address for NodeId")
}

var myId NodeId

func InitMyId() {
	if myId != 0 {
		panic("Cannot init my id multiple times")
	}
	myId = getIpAddress()
}

func GetMyId() NodeId {
	if myId == 0 {
		panic("My id has not been initialized")
	}
	return myId
}
