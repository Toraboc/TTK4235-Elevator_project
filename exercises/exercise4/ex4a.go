package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Data struct {
	mu    sync.Mutex
	Value int
	Stamp time.Time
}

func main() {
	start := time.Now()
	last := start
	d := Data{Value: 0, Stamp: time.Now().Add(-2 * time.Second)}

	go receiveUDP(&d)
	time.Sleep(100 * time.Millisecond)

	for {
		if time.Since(d.Stamp) >= 1*time.Second {
			fmt.Println("No recent updates received via UDP.")
			go sendUDP(&d)
			time.Sleep(100 * time.Millisecond)
			args := strings.Join(append([]string{os.Args[0]}, os.Args[1:]...), " ")
			cmd := exec.Command("gnome-terminal", "--", "bash", "-lc", args)
			_ = cmd.Start()
			break
		}
	}

	for {

		if time.Since(last) >= time.Second {
			d.Value++
			fmt.Printf("Time elapsed: %v seconds\n", d.Value)
			last = time.Now()
			if time.Since(start) >= time.Duration(rand.Intn(30))*time.Second {
				break
			}
		}
	}
}

func sendUDP(in *Data) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		in.mu.Lock()
		d := Data{Value: in.Value, Stamp: time.Now()}
		in.mu.Unlock()
		b, err := json.Marshal(d)
		if err != nil {
			continue
		}
		_, _ = conn.Write(b)
	}
}

func receiveUDP(out *Data) {
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	if err != nil {
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		var d Data
		if err := json.Unmarshal(buf[:n], &d); err != nil {
			continue
		}
		out.mu.Lock()
		out.Value = d.Value
		out.Stamp = d.Stamp
		out.mu.Unlock()
		fmt.Printf("Received: %d at %v\n", d.Value, d.Stamp)
	}
}
