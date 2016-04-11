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

			fmt.Println("Cost: ", cost, "for elev: ", elevatorDriver.ConnectedElevs[elev].IP)
	}
	return target
}

func getOwnCost(pos int, order elevatorDriver.Order) int{ 
	cost:= 0

	//elevator already at floor
	
	if elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor == order.Floor && elevatorDriver.ConnectedElevs[pos].Info.Dir == 0{ 
		cost = 0
		return cost
	}

	 //elevator already has orders at floor
	for button := 0; button < elevatorDriver.N_BUTTONS; button ++{ 
		if elevatorDriver.ConnectedElevs[pos].OwnQueue[order.Floor][button] == 1{
			cost = 1
			return cost
		}	
	}

	//higher cost for more orders and further away from ordered floor
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
					if elevatorDriver.ConnectedElevs[pos].OwnQueue[floor][button] == 1{
						cost += 1
						cost += int(math.Abs(float64(floor - elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor)))
					}
			}
	}



	return cost

}