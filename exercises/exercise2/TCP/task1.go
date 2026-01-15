package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("tcp", "10.100.23.11:33546") // delimiter port

	conn.(*net.TCPConn).SetNoDelay(true)

	go func(conn net.Conn) {
		defer conn.Close()

		reader := bufio.NewReader(conn)
		for {
			// Read messages until null
			line, err := reader.ReadBytes(0) // read until '\0'
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Received: \"%s\"\n", string(line[:len(line)-1])) // strip null
			}
		}
	}(conn)

	go func() {
		localPort := "20008"
		ln, err := net.Listen("tcp", "10.100.23.18:"+localPort)
		if err != nil {
			panic(err)
		}

		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		defer ln.Close()

		reader := bufio.NewReader(conn)
		for {
			// Read messages until null
			line, err := reader.ReadBytes(0) // read until '\0'
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Printf("Received: \"%s\"\n", string(line[:len(line)-1])) // strip null
			}
		}
	}()

	time.Sleep(1 * time.Second)
	// Send a message with null terminator
	msg := "Hello\x00"
	conn.Write([]byte(msg))

	for {
		time.Sleep(1 * time.Second)
		localIP := "10.100.23.18"
		localPort := "20008"
		msg = fmt.Sprintf("Connect to: %s:%s\x00", localIP, localPort)
		conn.Write([]byte(msg))
		time.Sleep(1 * time.Second)

	}
	time.Sleep(10 * time.Second)
}

/*
func main() {
	tcpSender := &net.TCPAddr{
		IP: net.ParseIP("10.100.23.11"),
		Port: 34933,
	}
	tcpListener := &net.TCPAddr{
		IP: net.ParseIP("0.0.0.0"),
		Port: 34933,
	}

	go func() {
		//listen and print
		conn, err := net.ListenTCP("tcp",tcpListener);
		if err != nil {
			fmt.Println("Error creating TCP socket")
		}
		defer conn.Close()

		conn.AcceptTCP()

	}()

}*/
