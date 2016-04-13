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
		fmt.Println("Calculating cost for: ", elevatorDriver.ConnectedElevs[elev].IP)
		cost := getOwnCost(elev, order)
		fmt.Println("Done calculating cost for: ", elevatorDriver.ConnectedElevs[elev].IP)
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
	info := queueDriver.GetInfoIP(elevatorDriver.ConnectedElevs[pos].IP)
	dir := info.Dir
	currentFloor := info.CurrentFloor
	fmt.Println("CurrentFloor: ", currentFloor)
	fmt.Println("Dir: ", dir)

	if currentFloor == order.Floor && (dir == 0) {
		//elevator already at floor
		cost = 0
		fmt.Println("At floor")
		return cost
	}

	for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
		//elevator already has orders at floor
		if elevatorDriver.ConnectedElevs[pos].CostQueue[order.Floor][button] == 1 {
			cost = 1
			fmt.Println("Order at floor")
			return cost
		}
	}
	i := 0
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			//fmt.Print(elevatorDriver.ConnectedElevs[pos].CostQueue[floor][button])
			if elevatorDriver.ConnectedElevs[pos].CostQueue[floor][button] == 1 {
				cost += 2 //higher cost for more orders
				i += 2

			}
		}
		//fmt.Println()
	}
	fmt.Println("Already has orders, cost: ", i)
	//fmt.Println()

	cost += 4 * int(math.Abs(float64(order.Floor-currentFloor))) //adds cost for distance from floor
	fmt.Println("Cost from different floor: ", 2*int(math.Abs(float64(order.Floor-currentFloor))))

	if dir == 1 && order.ButtonType == 1 { //elevator going up order going down
		cost += 3
		fmt.Println("Cost, dir up")
	} else if dir == -1 && order.ButtonType == 0 { //elevator going down order going up
		cost += 3
		fmt.Println("Cost dir down")
	}

	return cost

}
