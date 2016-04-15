package main

import (
	"./elevatorDriver"
	"./manager"
	"./network"
	"./queueDriver"
	"./userInterfaceDriver"
	"fmt"
	"os"
	"time"
)

var chButtonPressed = make(chan elevatorDriver.Order)
var chGetFloor = make(chan int)

var chToNetwork = make(chan network.Message, 100)
var chFromNetwork = make(chan network.Message, 100)

func main() {

	elevatorDriver.ElevInit()

	/*go func() {
		for {
			time.Sleep(25 * time.Millisecond)
			if elevatorDriver.ElevGetStopButton() {
				panic("Stop button pressed")
			}
		}
	}()*/

	if _, err := os.Open(elevatorDriver.QUEUE); err == nil {
		queueDriver.FileRead(elevatorDriver.QUEUE)
	} else {
		time.Sleep(time.Millisecond)
		if _, err := os.Create(elevatorDriver.QUEUE); err != nil {
			fmt.Println("Error, file not read")
		}
	}

	queueDriver.QueueInit()

	go userInterfaceDriver.NewOrder(chButtonPressed)
	go userInterfaceDriver.FloorTracker(chGetFloor)
	go manager.ChannelHandler(chButtonPressed, chGetFloor, chFromNetwork, chToNetwork)
	go network.NetworkHandler(chToNetwork, chFromNetwork)

	select {}
}
