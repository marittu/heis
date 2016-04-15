package manager

import (
	"../costManager"
	"../elevatorDriver"
	"../network"
	"../queueDriver"
	"fmt"
	"net"
	"strings"
	"time"
)

func ChannelHandler(chButtonPressed chan elevatorDriver.Order, chGetFloor chan int, chFromNetwork chan network.Message, chToNetwork chan network.Message) {
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(), "/")[0]
	timer := time.NewTimer(0)

	for {
		select {
		case order := <-chButtonPressed: //button pressed
			if order.ButtonType == 2 { //BUTTON_INTERNAL

				queueDriver.AddOrder(order)
				if elevatorDriver.Info.State == elevatorDriver.Idle{
					queueDriver.GetNextOrder(SelfIP, chToNetwork, timer)	
				}

				//Sending internal order to be added to the elevators CostQueue
				var temp elevatorDriver.ElevInfo
				temp.Dir = queueDriver.GetDir()
				temp.CurrentFloor = queueDriver.GetCurrentFloor()

				var msg network.Message
				msg.Order = order
				msg.Info = temp
				msg.FromIP = SelfIP
				msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
				msg.MessageId = network.NewInternalOrder

				chToNetwork <- msg

			} else { //External order
				//Sending external order to be added to the master queue and find target elevator
				var msg network.Message
				msg.Order = order
				msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
				msg.FromIP = SelfIP
				msg.MessageId = network.NewOrder

				chToNetwork <- msg
			}

		case floor := <-chGetFloor:
			fmt.Println("State: ", elevatorDriver.Info.State)
			queueDriver.PassingFloor(floor, SelfIP, chToNetwork, timer)

		case <-timer.C:
			timer.Stop()
			elevatorDriver.ElevSetDoorOpenLamp(0)
			elevatorDriver.Info.State = elevatorDriver.Idle
			if elevatorDriver.Info.State == elevatorDriver.Idle{
					queueDriver.GetNextOrder(SelfIP, chToNetwork, timer)	
			}

		case message := <-chFromNetwork:

			switch message.MessageId {

			case network.NewOrder:
				//Calculates the target elevator to take the order
				queueDriver.AddOrderMasterQueue(message.Order)
				if SelfIP == message.ToIP { //if master
					target := costManager.GetTargetElevator(message.Order)
					fmt.Println("target", target)

					//Sends order to target elevator for it to be added to CostQueue and its own queue
					var msg network.Message
					msg.Order = message.Order
					msg.ToIP = target
					msg.FromIP = SelfIP //master
					msg.MessageId = network.OrderFromMaster

					chToNetwork <- msg
				}

			case network.OrderFromMaster:
				if SelfIP == message.ToIP { //if master
					//Elevator recieved external order from master
					queueDriver.AddOrder(message.Order)
					if elevatorDriver.Info.State == elevatorDriver.Idle{
						queueDriver.GetNextOrder(SelfIP, chToNetwork, timer)	
					}
				}

			case network.Ack:
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS-1; button++ {
					elevatorDriver.ElevSetButtonLamp(message.Info.CurrentFloor, button, 0)
					queueDriver.MasterQueue[message.Info.CurrentFloor][button] = 0
				}

			case network.Removed:
				//Elevator that is removed from network takes all orders in masterQueue
				if len(elevatorDriver.ConnectedElevs) == 1{
					for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
						for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS-1; button++ {
							queueDriver.Queue[floor][button]= queueDriver.MasterQueue[floor][button]
						}
					}
				}

				//If removed elevator has unfinished orders - master handles them
				var added bool
				for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
					for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS-1; button++ {
						for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
							if queueDriver.MasterQueue[floor][button] == elevatorDriver.ConnectedElevs[elev].CostQueue[floor][button] {
								added = true
								break
							} else if queueDriver.MasterQueue[floor][button] != elevatorDriver.ConnectedElevs[elev].CostQueue[floor][button] {
								added = false
							}
						}
						if added == false {
							if SelfIP == elevatorDriver.ConnectedElevs[0].Master {
								order := elevatorDriver.Order{Floor: floor, ButtonType: button}
								queueDriver.AddOrder(order)
								fmt.Println("Adding orders from removed elevator: ", order)
							}
						}
					}
				}
				if elevatorDriver.Info.State == elevatorDriver.Idle{
					queueDriver.GetNextOrder(SelfIP, chToNetwork, timer)	
				}
			}
		}
	}
}
