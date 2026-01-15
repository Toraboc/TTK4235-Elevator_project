package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP("10.100.23.11"),
		Port: 20008, // server listening port
	}

	localAddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 20008, // local listening port
	}

	// Create a single UDP socket for sending and receiving
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatal("Error creating UDP socket:", err)
	}
	defer conn.Close()

	fmt.Printf("Listening on UDP port %d\n", localAddr.Port)

	buffer := make([]byte, 1024)
	i:=0
	// Goroutine: send messages every 5 seconds
	go func() {
		for {
			message := []byte("Hello server!")
			n, err := conn.WriteToUDP(message, serverAddr)
			if err != nil {
				fmt.Println("Error sending message:", err)
			} else {
				fmt.Printf("Sent %d bytes to %s\n", n, serverAddr.String())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Main loop: listen for any incoming messages
	for {
		i++
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				continue
			}
			fmt.Println("Error reading:", err)
			continue
		}

		// Only accept messages from the server IP
		if remoteAddr.IP.String() != serverAddr.IP.String() {
			continue
		}

		fmt.Printf("Received messag no. %d: %s from %s\n", i,string(buffer[:n]), remoteAddr)
	}
}
