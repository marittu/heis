package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	//"../network"
	//"fmt"
	//"time"
)



func ChannelHandler(chButtonPressed chan elevatorDriver.Button, chGetFloor chan int){
	//selectMaster()
	for{ 
		select{
		case order := <- chButtonPressed: //add case for internal order or external order
			queueDriver.AddOrder(order)
			queueDriver.GetDirection()
			break
		case floor := <- chGetFloor:
			queueDriver.PassingFloor(floor)
			break
		}
	}
}




