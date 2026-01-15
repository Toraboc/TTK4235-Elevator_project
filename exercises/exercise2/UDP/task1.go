package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

func main() {
	addr := net.UDPAddr{
		Port: 30000,
		IP:   net.ParseIP("10.100.23.155"),
	}
	if addr.IP == nil {
		log.Fatal("invalid IP address")
	}
	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		errors.New("cannot connect")
	}

	defer conn.Close()

	fmt.Println("LISTENING ON UDP PORT 30000")

	buffer := make([]byte, 1024)

	for {

		n, remoteAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading")
			continue
		}

		fmt.Printf("recieved message: %s from  %s\n", string(buffer[:n]), remoteAddr)
	}

}
