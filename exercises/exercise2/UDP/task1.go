package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	// Remote server to send messages to
	sendAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.100.23.155"), // server IP
		Port: 20008,                        // server port
	}

	// Local port to listen on
	listenerAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"), // listen on all interfaces
		Port: 30000,                  // local port
	}

	// Create sending socket (local port auto-assigned)
	sendConn, err := net.DialUDP("udp", nil, sendAddr)
	if err != nil {
		log.Fatal("Error dialing server:", err)
	}
	defer sendConn.Close()

	// Create listening socket
	listenConn, err := net.ListenUDP("udp", listenerAddr)
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listenConn.Close()

	fmt.Printf("Listening on UDP port %d\n", listenerAddr.Port)

	allowedIP := "10.100.23.155"
	buffer := make([]byte, 1024)

	// Goroutine for sending messages periodically
	go func() {
		for {
			message := []byte("Hello server!")
			n, err := sendConn.Write(message)
			if err != nil {
				fmt.Println("Error sending message:", err)
			} else {
				fmt.Printf("Sent %d bytes to %s\n", n, sendAddr.String())
			}
			time.Sleep(1 * time.Second) // send every second
		}
	}()

	// Main goroutine listens for incoming messages
	for {
		// Optional: set a read timeout to avoid blocking forever
		listenConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, remoteAddr, err := listenConn.ReadFromUDP(buffer)
		if err != nil {
			// Ignore timeout errors
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			fmt.Println("Error reading:", err)
			continue
		}

		// Only accept messages from the allowed IP
		if remoteAddr.IP.String() != allowedIP {
			continue
		}

		fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), remoteAddr)
	}
}
