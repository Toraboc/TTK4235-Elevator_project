package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"sync"
	"time"
	"encoding/binary"
	"io"
)

var number int32
var mu sync.Mutex

const filename = "exercise4.go"
const socketFile = "/home/student/Documents/Sanntid55/exercises/exercise4/heartbeat.sock"

func main() {
	backupState()

	// Spawn backup
	spawnNewProcess()

	mainState()
}

func spawnNewProcess() {
	os := runtime.GOOS
	switch os {
	case "darwin":
		exec.Command("osascript", "-e", `tell app "Terminal" to do script "go run `+filename+`"`).Run()
	case "linux":
		exec.Command("gnome-terminal", "--", "go", "run", filename).Run()
	default:
		panic("Operating system not supported")
	}
}

func backupState() {
	fmt.Println("Entering Backup")
	listener, err := net.Listen("unix", socketFile)
	if err != nil {
		fmt.Println(err)
		panic("Failed to listen to file socket")
	}
	defer listener.Close()
	unixListener := listener.(*net.UnixListener)

	for {
		unixListener.SetDeadline(time.Now().Add(1000 * time.Millisecond))
		conn, err := unixListener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				// Timeout, the other process is probably dead
				return
			}
			fmt.Println("Failed to receive a message")
		}

		buf := make([]byte, 4)

		for {
			_, err2 := conn.Read(buf)
			if err2 != nil {
				if err2 == io.EOF {
					fmt.Println("Client disconnected")
					// Probably a crash
					return
				} else {
					fmt.Println("Read error:", err2)
				}
				break
			}

			val := int32(binary.BigEndian.Uint32(buf))
			mu.Lock()
			number = val
			mu.Unlock()
		}
	}
}

func mainState() {
	fmt.Println("Taking over as PRIMARY")
	conn, err := getConnection()
	if err != nil {
		fmt.Println(err)
		panic("Failed to open socket file")
	}

	go printLoop()
	go heartBeatLoop(conn)

	for {}
}

func getConnection() (net.Conn, error) {
	deadline := time.Now().Add(500 * time.Millisecond)

    for {
        conn, err := net.Dial("unix", socketFile)
        if err == nil {
            return conn, nil
        }

        if time.Now().After(deadline) {
            return nil, err
        }

        time.Sleep(50 * time.Millisecond)
    }
}

func printLoop() {
	for {
		mu.Lock()
		number++
		mu.Unlock()
		fmt.Println(number)
		time.Sleep(1 * time.Second)
	}
}

func heartBeatLoop(conn net.Conn) {
	for {
		mu.Lock()
		currentNumber := number
		mu.Unlock()
		// Send heartbeat

		buf := make([]byte, 4) // 4 bytes for int32
		binary.BigEndian.PutUint32(buf, uint32(currentNumber))

		_, err := conn.Write(buf)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
