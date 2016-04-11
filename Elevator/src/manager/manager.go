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

		/*for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
			fmt.Println(elevatorDriver.ConnectedElevs[elev])
		}*/
		select{
		case order := <- chButtonPressed: //button pressed


			if order.ButtonType == 2{ //BUTTON_INTERNAL 
				
				fmt.Println(SelfIP)

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
				
			}else{ //External order

				
				fmt.Println("Order recieved")
				var msg network.Message
				msg.Order = order
				msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
				msg.FromIP = SelfIP
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
				if SelfIP == message.ToIP{ //if master
					target := costManager.GetTargetElevator(message.Order)
					fmt.Println("target", target)	
				}

			/*case network.NewInternalOrder:
				for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
					if elevatorDriver.ConnectedElevs[elev].IP == message.FromIP{
					}
				}
			*/

								
			}
		}
	}
}



