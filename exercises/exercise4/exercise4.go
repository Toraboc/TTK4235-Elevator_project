package main

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

var number int
var mu sync.Mutex

const filename = "exercise4.go"

func main() {
	backupState()

	fmt.Println("Taking over...")

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
	quit := make(chan bool)
	lastHB := time.Now()
	var HBMu sync.Mutex

	go func() {
		listener, err := net.Listen("unix", "./heartbeat.sock")
		if err != nil {
			fmt.Println(err)
			panic("Failed to listen to file socket")
		}
		defer listener.Close()
		tcpListener := listener.(*net.TCPListener)

		for {
			select {
			case <-quit:
				return
			default:
				tcpListener.SetDeadline(time.Now().Add(300 * time.Millisecond))
				conn, err := tcpListener.Accept()
				if err != nil {
					if ne, ok := err.(net.Error); ok && ne.Timeout() {
						// Accept timed out
						continue // or break / handle timeout
					}
					fmt.Println("Failed to receive a message")
				}

				buf := make([]byte, 16)
				nr, err := conn.Read(buf)
				if err != nil {
					return
				}

				data := buf[0:nr]
				mu.Lock()
				number = int(data[0])
				mu.Unlock()
				
				HBMu.Lock()
				lastHB = time.Now()
				HBMu.Unlock()

			}
		}
	}()

	for {
		time.Sleep(300 * time.Millisecond)

		HBMu.Lock()
		timeDiff := time.Since(lastHB)
		HBMu.Unlock()

		if (timeDiff > 300 * time.Millisecond) {
			quit <- true
			break
		}
	}
}

func mainState() {
	go printLoop()
	go heartBeatLoop()

	time.Sleep(5 * time.Second)
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

func heartBeatLoop() {
	for {
		time.Sleep(100 * time.Millisecond)
		mu.Lock()
		currentNumer := number
		mu.Unlock()
		// Send heartbear
		fmt.Println("HB:", currentNumer)
	}
}
