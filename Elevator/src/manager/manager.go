package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	"../network"
	//"fmt"
	//"time"
)



func ChannelHandler(chButtonPressed chan elevatorDriver.Button, chGetFloor chan int, chFromNetwork chan network.Message, chToNetwork chan network.Message){
	elevator := network.ElevManagerInit()
	for{ 
		select{
		case order := <- chButtonPressed: //button pressed
			if order.ButtonType == 2{ //BUTTON_INTERNAL
				//fmt.Println("New internalOrder")
				queueDriver.AddOrder(order)
				queueDriver.GetDirection()
	
			}else{ //External order
				//fmt.Println("New external order")
				var msg network.Message
				msg.Order = order
				msg.ToIP = elevator.Master
				msg.FromIP = elevator.SelfIP
				msg.MessageId = network.NewOrder

				chToNetwork <- msg
			}
			
			break
		case floor := <- chGetFloor:
			queueDriver.PassingFloor(floor)
			break
		case message := <-chFromNetwork:
			//fmt.Println("Message id: ", message.MessageId)
			switch(message.MessageId){

			case 2: //New order
				queueDriver.AddOrderMasterQueue(message.Order)
				//queueDriver.PrintQueue()

			}
		}
	}
}




