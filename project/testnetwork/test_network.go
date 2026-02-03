package testnetwork

import (
	"bytes"
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	// port is the UDP port used for both listening and broadcasting.
	port = 42069
	// broadcast is the IPv4 broadcast address and port used for discovery.
	broadcast = "255.255.255.255:42069"
	// sendHz is the broadcast frequency in Hz.
	sendHz = 200
	// printEvery is the interval used for peer list logging.
	printEvery = time.Second
)

// nodeSet tracks last-seen timestamps for peers by IP address.
type nodeSet struct {
	mu   sync.Mutex
	last map[string]time.Time
}

// newNodeSet creates an initialized nodeSet.
func newNodeSet() *nodeSet {
	return &nodeSet{last: make(map[string]time.Time)}
}

// seen records that the given IP was observed now.
func (s *nodeSet) seen(ip string) {
	s.mu.Lock()
	s.last[ip] = time.Now()
	s.mu.Unlock()
}

// list returns the sorted list of active peer IPs and prunes stale entries.
func (s *nodeSet) list() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	// drop stale nodes
	now := time.Now()
	for ip, t := range s.last {
		if now.Sub(t) > time.Second/2 {
			delete(s.last, ip)
		}
	}
	ips := make([]string, 0, len(s.last))
	for ip := range s.last {
		ips = append(ips, ip)
	}
	sort.Strings(ips)
	return ips
}

// Message represents a payload received from a peer.
type Message struct {
	// FromIP is the sender's IP address.
	FromIP string
	// Payload is the received data payload.
	Payload []byte
}

// Start launches UDP broadcast sender and receiver.
// It broadcasts the local IP plus the provided payload at sendHz and prints
// discovered peers at a regular interval. This call blocks and only returns
// on read error.
//
// Use StartWithHandler to receive payloads from peers.
func Start(data []byte) error {
	return StartWithHandler(data, nil)
}

// StartWithHandler is like Start, but invokes onMessage for each received
// payload. The handler is called synchronously on the receiver loop, so keep
// it fast or offload work to another goroutine.
func StartWithHandler(data []byte, onMessage func(Message)) error {
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		return err
	}
	defer conn.Close()

	peers := newNodeSet()

	go func() {
		t := time.NewTicker(time.Second / sendHz)
		defer t.Stop()

		bcastAddr, _ := net.ResolveUDPAddr("udp4", broadcast)
		for range t.C {
			ip := localIP()
			if ip == "" {
				continue
			}
			msg := append([]byte(ip+" "), data...)
			_, _ = conn.WriteToUDP(msg, bcastAddr)
		}
	}()

	go func() {
		t := time.NewTicker(printEvery / 10)
		defer t.Stop()
		for range t.C {
			ips := peers.list()
			fmt.Printf("peers: %d %v\n", len(ips), ips)
		}
	}()

	buf := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		ip := addr.IP.String()
		peers.seen(ip)

		if onMessage != nil {
			parts := bytes.SplitN(buf[:n], []byte(" "), 2)
			fromIP := ip
			payload := buf[:n]
			if len(parts) == 2 {
				fromIP = string(parts[0])
				payload = parts[1]
			}
			payloadCopy := make([]byte, len(payload))
			copy(payloadCopy, payload)
			onMessage(Message{FromIP: fromIP, Payload: payloadCopy})
		}
	}
}

// localIP returns the first non-loopback IPv4 address of an active interface.
// If no suitable address is found, it returns an empty string.
func localIP() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			ip4 := ipnet.IP.To4()
			if ip4 != nil {
				return ip4.String()
			}
		}
	}
	return ""
}
