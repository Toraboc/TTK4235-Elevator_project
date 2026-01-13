// Use `go run foo.go` to run your program

package main

import (
	"fmt"
	"runtime"
)

type request struct {
	action string
	reply  chan int
}

func server(req chan request) {
	i := 0
	for {
		select {
		case r := <-req:
			switch r.action {
			case "inc":
				i++
			case "dec":
				i--
			case "get":
				r.reply <- i
			}
		}
	}
}

func incrementing(req chan request, done chan bool) {
	for i := 0; i < 10000; i++ {
		req <- request{action: "inc"}
	}
	done <- true
}

func decrementing(req chan request, done chan bool) {
	for i := 0; i < 10000; i++ {
		req <- request{action: "dec"}
	}
	done <- true
}

func main() {
	// What does GOMAXPROCS do? What happens if you set it to 1?
	runtime.GOMAXPROCS(3)

	req := make(chan request)
	done := make(chan bool)

	go server(req)

	go incrementing(req, done)
	go decrementing(req, done)

	<-done
	<-done

	reply := make(chan int)
	req <- request{action: "get", reply: reply}
	final := <-reply

	fmt.Println("The magic number is:", final)
}
