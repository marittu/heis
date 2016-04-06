package main

import (
	"./elevatorDriver"
	"./userInterfaceDriver"
	"./queueDriver"
	"./manager"
	//"fmt"
	//"time"
	"./network"
	
)

var chButtonPressed = make(chan elevatorDriver.Button)
var chGetFloor = make(chan int)
var chToNetwork = make(chan Message, 100)
var chFromNetwork = make(chan Message, 100)

func main() {

	
	//queueDriver.QueueInit()
	//elevatorDriver.ElevInit()

	
	//go userInterfaceDriver.NewOrder(chButtonPressed)
	//go userInterfaceDriver.FloorTracker(chGetFloor)
	//go manager.ChannelHandler(chButtonPressed, chGetFloor)
	go network.manager(chToNetwork, chFromNetwork)	
	

	for{}
}
