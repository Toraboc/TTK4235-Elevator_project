package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"  // <-- add this line
	"time"
)


func main() {
	role := "primary"
	if len(os.Args) > 1 {
		role = os.Args[1]
	}

	if role == "primary" {
		runPrimary()
	} else {
		runBackup()
	}
}

func runPrimary() {
	fmt.Println("PRIMARY:", os.Getpid())

	// Start the backup in a new terminal
	cmd := exec.Command("cmd", "/C", "start", "cmd", "/K", os.Args[0], "backup")
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting backup:", err)
	}

	i := 1

	// Read heartbeat file to continue counting
	data, err := os.ReadFile("heartbeat.txt")
	if err == nil {
		if n, err2 := strconv.Atoi(string(data)); err2 == nil {
			i = n + 1
		}
	}

	for {
		fmt.Println(i)
		os.WriteFile("heartbeat.txt", []byte(fmt.Sprint(i)), 0644)
		i++
		time.Sleep(1 * time.Second)
	}
}


func runBackup() {
	fmt.Println("BACKUP PID:", os.Getpid())

	timeout := 3 * time.Second
	var lastSeen time.Time

	for {

		info, err := os.Stat("heartbeat.txt")
		if err == nil {
			lastSeen = info.ModTime()
		} else {
			fmt.Println("Backup: heartbeat file missing")
		}

		if !lastSeen.IsZero() && time.Since(lastSeen) > timeout {
			fmt.Println("Backup: TIMEOUT — promoting!")

			cmd := exec.Command(os.Args[0], "primary")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Start()
			fmt.Println("Backup: start result:", err)

			os.Exit(0)
		}

		time.Sleep(1 * time.Second)
	}
}

