package main

import (
	"fmt"
	"project/drivers/network"
)

func main2() {
	data := []byte("Hello, peers!")
	handleMessage := func(msg network.Message) {
		fmt.Printf("from %s: %s\n", msg.FromIP, string(msg.Payload))
	}
	err := network.StartWithHandler(data, handleMessage)
	if err != nil {
		panic(err)
	}
}
