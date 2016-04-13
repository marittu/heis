package costManager

import (
	"../elevatorDriver"
	"../queueDriver"
	"fmt"
	"math"
)

func GetTargetElevator(order elevatorDriver.Order) string {
	target := elevatorDriver.ConnectedElevs[0].IP
	min := 100
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		cost := getOwnCost(elev, order)
		if cost < min {
			min = cost
			target = elevatorDriver.ConnectedElevs[elev].IP

		}

		fmt.Println("Cost: ", cost, "for elev: ", elevatorDriver.ConnectedElevs[elev].IP) //fjern før levering
	}
	return target
}

func getOwnCost(pos int, order elevatorDriver.Order) int {

	cost := 0
	cur := queueDriver.GetCurrentFloorIP(elevatorDriver.ConnectedElevs[pos].IP)
	//fmt.Println("CurrentFloor ", cur, "for elev ", elevatorDriver.ConnectedElevs[pos].IP)

	if (cur == order.Floor) && (elevatorDriver.ConnectedElevs[pos].Info.Dir == 0) {
		//elevator already at floor

		cost = 0
	}

	for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
		//elevator already has orders at floor
		if elevatorDriver.ConnectedElevs[pos].CostQueue[order.Floor][button] == 1 {
			cost = 1
		}
	}

	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			fmt.Print(elevatorDriver.ConnectedElevs[pos].CostQueue[floor][button])
			if elevatorDriver.ConnectedElevs[pos].CostQueue[floor][button] == 1 {
				cost += 5 //higher cost for more orders
				//fmt.Println("Already has orders")
			}
		}
		fmt.Println()
	}
	fmt.Println()

	cost += 2 * int(math.Abs(float64(order.Floor-cur))) //adds cost for distance from floor
	//fmt.Println("Cost from different floor: ", int(math.Abs(float64(order.Floor-cur))))

	//funker dir?
	if elevatorDriver.ConnectedElevs[pos].Info.Dir == 1 && order.ButtonType == 1 { //elevator going up order going down
		cost += 3
	} else if elevatorDriver.ConnectedElevs[pos].Info.Dir == -1 && order.ButtonType == 0 { //elevator going down order going up
		cost += 3
	}

	return cost

}
