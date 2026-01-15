package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Local address (where we listen for the reply)
	localAddr := &net.UDPAddr{
		IP:   net.IPv4zero, // listen on all local interfaces
		Port: 20008,
	}

	// Remote server address
	remoteAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.100.23.155"),
		Port: 30000,
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		log.Fatal("DialUDP failed:", err)
	}
	defer conn.Close()

	// Send message
	message := []byte("hello from client")
	_, err = conn.Write(message)
	if err != nil {
		log.Fatal("Write failed:", err)
	}

	fmt.Println("Message sent, waiting for reply...")

	// Optional: avoid blocking forever
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Read reply
	buffer := make([]byte, 1024)
	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal("Read failed:", err)
	}

	fmt.Printf("Received reply from %s: %s\n", addr, string(buffer[:n]))
}
