package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	"../network"
	//"../costManager"
	"fmt"
	"net"
	"strings"
	//"time"
)



func ChannelHandler(chButtonPressed chan elevatorDriver.Order, chGetFloor chan int, chFromNetwork chan network.Message, chToNetwork chan network.Message){
	//elevator := network.GetElevManager()
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(),"/")[0]
	for{ 
		select{
		case order := <- chButtonPressed: //button pressed


			if order.ButtonType == 2{ //BUTTON_INTERNAL
				queueDriver.AddOrder(order, SelfIP)
				queueDriver.GetDirection(SelfIP)
				
			}else{ //External order

				
				fmt.Println("Order recieved")
				var msg network.Message
				msg.Order = order
				msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
				//msg.FromIP = SelfIP
				msg.MessageId = network.NewOrder

				chToNetwork <- msg
				
			}
			
			break
		case floor := <- chGetFloor:
			queueDriver.PassingFloor(floor, SelfIP)
			break
		case message := <-chFromNetwork:
			
			switch(message.MessageId){

			case network.NewOrder:
				queueDriver.AddOrderMasterQueue(message.Order)

				/*
			case 4: //Find target
				if SelfIP == message.ToIP{
					target := network.GetMinCost()
					fmt.Println("Target: ", target)
					var msg network.Message
					msg.ToIP = target
					msg.FromIP = elevator.Master
					msg.MessageId = network.MasterDistributesOrder
					msg.Order = message.Order
					//fmt.Println("sending target")
					chToNetwork <- msg //fyller opp kanalen - kommer ikke videre
					//fmt.Println("send done")
	
				}
				*/
				//network.AppendCost(message.FromIP, message.Cost)
					//fmt.Println("Sending")
				//fmt.Println("Elevator: ", message.FromIP, " cost: ", message.Cost)

				//queueDriver.AddOrderMasterQueue(message.Order)
				
				/*for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{

					
				}*/	

				//target := costManager.GetTargetElev()
				
			
			/*
			case 5: //Target adds order from master
				//fmt.Println("in add order case")
				if SelfIP == message.ToIP{
					//fmt.Println("Adding order")
					queueDriver.AddOrder(message.Order)
					queueDriver.GetDirection()
				}
			*/	
			}
		}
	}
}



