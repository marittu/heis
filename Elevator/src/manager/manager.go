package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	"../network"
	"../costManager"
	"fmt"
	//"net"
	//"strings"
	//"time"
)



func ChannelHandler(chButtonPressed chan elevatorDriver.Order, chGetFloor chan int, chFromNetwork chan network.Message, chToNetwork chan network.Message){
	elevator := network.GetElevManager()
	/*addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(),"/")[0]*/
	for{ 
		select{
		case order := <- chButtonPressed: //button pressed


			if order.ButtonType == 2{ //BUTTON_INTERNAL
				queueDriver.AddOrder(order)
				queueDriver.GetDirection()
	
			}else{ //External order
				//target := costManager.GetTargetElev(order, elevator.SelfIP)
				for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
					cost := costManager.GetOwnCost(order)

					//network.SendNetworkMessage(cost, elevator.SelfIP, elevator.Master, network.NewOrder, chToNetwork)
					var msg network.Message
					msg.Cost = cost
					msg.ToIP = elevator.Master
					msg.FromIP = elevatorDriver.ConnectedElevs[elev].IP
					msg.MessageId = network.NewOrder

					chToNetwork <- msg

					fmt.Println("Elevator: ", msg.FromIP, " cost: ", cost)
				}
			}
			
			break
		case floor := <- chGetFloor:
			queueDriver.PassingFloor(floor)
			break
		case message := <-chFromNetwork:
			
			switch(message.MessageId){

			case 2: //New order

				queueDriver.AddOrderMasterQueue(message.Order)
				network.AppendCost(message.FromIP, message.Cost)

				/*for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{

					
				}*/	

				//target := costManager.GetTargetElev()
				
			
			
			//case 3: //Master sends order to target elev
				
			}
		}
	}
}



