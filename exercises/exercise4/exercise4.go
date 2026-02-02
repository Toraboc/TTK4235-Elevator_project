package main

import (
	"fmt"
	// "github.com/syucream/posix_mq/src/posix_mq"
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
		exec.Command("osascript", "-e", `tell app "Terminal" to do script "go run ` + filename + `"`).Run()
	case "linux":
		exec.Command("gnome-terminal", "--", "go", "run", filename).Run()
	default:
		panic("Operating system not supported")
	}
}

func backupState() {
	fmt.Println("Backing up")
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
