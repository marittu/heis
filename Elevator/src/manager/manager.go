package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	"../network"
	"../costManager"
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
				
				queueDriver.AddOrder(order) // , SelfIP
				queueDriver.GetDirection(SelfIP)
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
				
				break	
			}else{ //External order
				
				queueDriver.AddOrderMasterQueue(order)
				//fmt.Println("Order recieved")
				var msg network.Message
				msg.Order = order
				msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
				msg.FromIP = SelfIP
				msg.MessageId = network.NewOrder
			
				chToNetwork <- msg
				
				
				break
			}
			
			break
		case floor := <- chGetFloor:
			queueDriver.PassingFloor(floor, SelfIP)
			break
		case message := <-chFromNetwork:
			
			switch(message.MessageId){

			case network.NewOrder:
				
				if SelfIP == message.ToIP{ //if master
					target := costManager.GetTargetElevator(message.Order)
					fmt.Println("target", target)
					var msg network.Message
					msg.Order = message.Order
					msg.ToIP = target
					msg.FromIP = SelfIP //master
					msg.MessageId = network.OrderFromMaster

					chToNetwork <- msg	
				}

			case network.OrderFromMaster:
				for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{

					if elevatorDriver.ConnectedElevs[elev].IP == message.ToIP{
						queueDriver.AddOrder(message.Order)
						queueDriver.GetDirection(SelfIP)

					}
				}
			

								
			}
		}
	}
}



