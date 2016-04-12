package main

import (
	"./elevatorDriver"
	"./manager"
	"./queueDriver"
	"./userInterfaceDriver"
	//"fmt"
	//"time"
	"./network"
)

var chButtonPressed = make(chan elevatorDriver.Order)
var chGetFloor = make(chan int)
var chToNetwork = make(chan network.Message, 100)
var chFromNetwork = make(chan network.Message, 100)

func main() {

	queueDriver.QueueInit()
	elevatorDriver.ElevInit()

	go userInterfaceDriver.NewOrder(chButtonPressed)
	go userInterfaceDriver.FloorTracker(chGetFloor)
	go manager.ChannelHandler(chButtonPressed, chGetFloor, chFromNetwork, chToNetwork)
	go network.NetworkHandler(chToNetwork, chFromNetwork)

	for {
	}

}
