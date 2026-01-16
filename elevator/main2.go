package main

import (
	"elevator/network"
)

func main() {
	data := []byte("Hello, peers!")
	err := network.Start(data)
	if err != nil {
		panic(err)
	}
}
