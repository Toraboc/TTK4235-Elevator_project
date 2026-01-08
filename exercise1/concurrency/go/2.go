// Use `go run foo.go` to run your program

package main

import (
    . "fmt"
    "runtime"
    "sync"
)

var i = 0
var mu sync.Mutex
var wg sync.WaitGroup

func incrementing() {
    //TODO: increment i 1000000 times
    defer wg.Done()
    for j := 0; j < 10000; j++ {
        mu.Lock()
        i++
        println(i)
        mu.Unlock()
    }
}

func decrementing() {
    //TODO: decrement i 1000000 times
    defer wg.Done()
    for j := 0; j < 10000; j++ {
        mu.Lock()
        i--
        println(i)
        mu.Unlock()
    }
}

func main() {
    // What does GOMAXPROCS do? What happens if you set it to 1?
    runtime.GOMAXPROCS(2)    
	
    wg.Add(2)
    // TODO: Spawn both functions as goroutines
    go incrementing()
    go decrementing()

    wg.Wait()

    // We have no direct way to wait for the completion of a goroutine (without additional synchronization of some sort)
    // We will do it properly with channels soon. For now: Sleep.
    Println("The magic number is:", i)
}