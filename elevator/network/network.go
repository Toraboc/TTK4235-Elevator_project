package network

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	port       = 42069
	broadcast  = "255.255.255.255:42069"
	sendHz     = 200
	printEvery = time.Second
)

type nodeSet struct {
	mu   sync.Mutex
	last map[string]time.Time
}

func newNodeSet() *nodeSet {
	return &nodeSet{last: make(map[string]time.Time)}
}

func (s *nodeSet) seen(ip string) {
	s.mu.Lock()
	s.last[ip] = time.Now()
	s.mu.Unlock()
}

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

// Start launches UDP broadcast sender and receiver.
// It broadcasts the local IP and data at 200 Hz and prints discovered peers.
func Start(data []byte) error {
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
		_ = n
		peers.seen(ip)
	}
}

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
