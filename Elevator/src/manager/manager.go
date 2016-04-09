package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	"../network"
	"../costManager"
	//"fmt"
	"net"
	"strings"
	//"time"
)



func ChannelHandler(chButtonPressed chan elevatorDriver.Order, chGetFloor chan int, chFromNetwork chan network.Message, chToNetwork chan network.Message){
	elevator := network.GetElevManager()
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(),"/")[0]
	for{ 
		select{
		case order := <- chButtonPressed: //button pressed


			if order.ButtonType == 2{ //BUTTON_INTERNAL
				queueDriver.AddOrder(order)
				queueDriver.GetDirection()
	
			}else{ //External order
				//target := costManager.GetTargetElev(order, elevator.SelfIP)
				cost := costManager.GetOwnCost(order)

				//network.SendNetworkMessage(cost, elevator.SelfIP, elevator.Master, network.NewOrder, chToNetwork)
				var msg network.Message
				msg.Cost = cost
				msg.ToIP = elevator.Master
				msg.FromIP = SelfIP
				msg.MessageId = network.NewOrder

				chToNetwork <- msg
			}
			
			break
		case floor := <- chGetFloor:
			queueDriver.PassingFloor(floor)
			break
		case message := <-chFromNetwork:
			
			switch(message.MessageId){

			case 2: //New order
				//fmt.Println("IP ", message.FromIP)
				//ownCost := costManager.GetOwnCost(message.Order)
				//fmt.Println("self ", SelfIP)
				network.AppendCost(message.FromIP)	
				//fmt.Println("Cost: ", message.Cost)
				/*for len(cost) < len(elevatorDriver.ConnectedElevs){
					
				}*/
			
				queueDriver.AddOrderMasterQueue(message.Order)
				//cost := costManager.GetOwnCost(message.Order) 
				//network.sendNetworkMessage(order, elevator.SelfIP, elevator.Master, network.Cost,  chToNetwork)
				//fmt.Println("Target: ", target)
				//queueDriver.PrintQueue()
			
			
			//case 3: //Master finds target elev
				//costManager.FindTarget()
			}
		}
	}
}



