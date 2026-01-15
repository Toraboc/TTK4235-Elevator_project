package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    // Define the UDP address to listen on
    addr, err := net.ResolveUDPAddr("udp", ":30000")
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        os.Exit(1)
    }

    // Create a UDP connection
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Error setting up UDP listener:", err)
        os.Exit(1)
    }
    defer conn.Close()

    fmt.Println("Listening on UDP port 30000...")

    // Buffer to store incoming data
    buffer := make([]byte, 1024)

    for {
        // Read data from the connection
        n, remoteAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error reading from UDP connection:", err)
            continue
        }

        // Print the received message and the sender's address
        fmt.Printf("Received message: %s from %s\n", string(buffer[:n]), remoteAddr)
    }
}