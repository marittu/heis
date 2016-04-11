package costManager

import(
	"../elevatorDriver"
	//"../network"
	//"../queueDriver"
	"fmt"
	"math"
)

func GetTargetElevator(order elevatorDriver.Order) string{
	target := elevatorDriver.ConnectedElevs[0].IP
	min := 100
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
		cost := getOwnCost(elev, order)
		if cost < min{
			min = cost
			target = elevatorDriver.ConnectedElevs[elev].IP

		}

			
	}
	return target
}

func getOwnCost(pos int, order elevatorDriver.Order) int{ 
	
	


	cost:= 0

	//elevator already at floor
	
	if (elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor == order.Floor) && (elevatorDriver.ConnectedElevs[pos].Info.Dir == 0){ 
		cost = 0
		//fmt.Println("cost 0")
		return cost
	}

	 //elevator already has orders at floor
	for button := 0; button < elevatorDriver.N_BUTTONS; button ++{ 
		if elevatorDriver.ConnectedElevs[pos].OwnQueue[order.Floor][button] == 1{
			cost = 1
			//fmt.Println("cost 1")
			return cost
		}	
	}

	//higher cost for more orders 
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
					if elevatorDriver.ConnectedElevs[pos].OwnQueue[floor][button] == 1{
						cost += 1
						//fmt.Println("cost more orders 1")
						
					}
			}
	}

	cost += int(math.Abs(float64(order.Floor - elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor)))
	//fmt.Println("Distance from floor ",int(math.Abs(float64(order.Floor - elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor))), "CurrentFloor: ", elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor )
	//fmt.Println("cost away from floor")
	if elevatorDriver.ConnectedElevs[pos].Info.Dir == 0{ //elevator at floor
		cost += 1 
		//fmt.Println("cost dir 0")
	}else if elevatorDriver.ConnectedElevs[pos].Info.Dir == 1 && order.ButtonType == 1{ //elevator going up order going down
		cost += 3
		//fmt.Println("cost dir 1")
	}else if elevatorDriver.ConnectedElevs[pos].Info.Dir == -1 && order.ButtonType == 0{ //elevator going down order going up
		cost += 3
		//fmt.Println("cost dir -1")
	}


	//fmt.Println("Cost: ", cost, "for elev: ", elevatorDriver.ConnectedElevs[pos].IP)
	return cost

}