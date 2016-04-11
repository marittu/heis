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
				queueDriver.GetDirection(SelfIP, chToNetwork)


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
				
			}else{ //External order
				
				queueDriver.AddOrderMasterQueue(order)
				var msg network.Message
				msg.Order = order
				msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
				msg.FromIP = SelfIP
				msg.MessageId = network.NewOrder
			
				chToNetwork <- msg
				
			}
			
			break
		
		case floor := <- chGetFloor:
			fmt.Println("Recieved from floor channel")
			queueDriver.PassingFloor(floor, SelfIP, chToNetwork)
			break
		
		/*case floor := <- chDoorOpen:
			fmt.Println("Kommer hit?")
			
			var temp elevatorDriver.ElevInfo
			temp.Dir = 0
			temp.CurrentFloor = floor
			var msg network.Message
			msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
			msg.Info = temp
			msg.MessageId = network.Ack
				
			chToNetwork <- msg
			fmt.Println("Stuck?")

			break*/


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
				if SelfIP == message.ToIP{ //if master
					fmt.Println("Order to: ", message.ToIP)
					queueDriver.AddOrder(message.Order)
					queueDriver.PrintQueue()
					queueDriver.GetDirection(SelfIP, chToNetwork)
					break
				}
				
			case network.Ack:
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++{
					elevatorDriver.ElevSetButtonLamp(message.Info.CurrentFloor,button,0) 	
				}		

								
			}
		}
	}
}



