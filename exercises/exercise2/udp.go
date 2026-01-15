package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// Channel to pass the server IP from listener to sender
	serverIPChan := make(chan string)

	// Start listener for server broadcast
	go listenForServerBroadcast(serverIPChan)

	// Wait for server IP from broadcast
	serverIP := <-serverIPChan
	fmt.Printf("Discovered server IP: %s\n", serverIP)

	// Start sending messages to the server
	sendMessages(serverIP)
}

// Listen for UDP broadcasts on port 30000 to discover server IP
func listenForServerBroadcast(serverIPChan chan<- string) {
	addr := net.UDPAddr{
		Port: 30000,
		IP:   net.IPv4zero,
	}

	conn, err := net.ListenUDP("udp4", &addr)
	if err != nil {
		fmt.Printf("Error listening for UDP: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Printf("Error reading UDP: %v\n", err)
			continue
		}
		message := strings.TrimSpace(string(buffer[:n]))
		fmt.Printf("Received broadcast from %s: %s\n", remoteAddr.IP, message)

		// Send server IP to channel
		serverIPChan <- remoteAddr.IP.String()
		return
	}
}

// Send messages to the server and listen for responses
func sendMessages(serverIP string) {
	serverAddr := fmt.Sprintf("%s:20008", serverIP)
	udpAddr, err := net.ResolveUDPAddr("udp4", serverAddr)
	if err != nil {
		fmt.Printf("Error resolving server address: %v\n", err)
		os.Exit(1)
	}

	// Create UDP connection for sending and receiving
	conn, err := net.DialUDP("udp4", nil, udpAddr)
	if err != nil {
		fmt.Printf("Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// Goroutine to listen for server responses
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				fmt.Printf("Error reading from UDP: %v\n", err)
				continue
			}
			fmt.Printf("Received from server: %s\n", string(buffer[:n]))
		}
	}()

	// Send messages periodically
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Hello server, message %d", i)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Printf("Error sending UDP message: %v\n", err)
		} else {
			fmt.Printf("Sent: %s\n", message)
		}
		time.Sleep(2 * time.Second)
	}
}
