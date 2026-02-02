package main

import (
	"fmt"
	"project/testnetwork"
)

func main2() {
	data := []byte("Hello, peers!")
	handleMessage := func(msg testnetwork.Message) {
		fmt.Printf("from %s: %s\n", msg.FromIP, string(msg.Payload))
	}
	err := testnetwork.StartWithHandler(data, handleMessage)
	if err != nil {
		panic(err)
	}
}
