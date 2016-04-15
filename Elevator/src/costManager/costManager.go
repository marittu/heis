package costManager

import (
	"../elevatorDriver"
	"../queueDriver"
	//"fmt"
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
	}
	return target
}

func getOwnCost(pos int, order elevatorDriver.Order) int {
	if elevatorDriver.ConnectedElevs[pos].Info.TimedOut == false{
		cost := 0
		info := queueDriver.GetInfoIP(elevatorDriver.ConnectedElevs[pos].IP)
		dir := info.Dir
		currentFloor := info.CurrentFloor
		if currentFloor == order.Floor && (dir == 0) {
			//elevator already at floor
			cost = 0
			return cost
		}

		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			//elevator already has orders at floor
			if elevatorDriver.ConnectedElevs[pos].CostQueue[order.Floor][button] == 1 {
				cost = 1
				return cost
			}
		}
		for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
			for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
				if elevatorDriver.ConnectedElevs[pos].CostQueue[floor][button] == 1 {
					cost += 2 //higher cost for more orders
				}
			}
		}

		cost += 4 * int(math.Abs(float64(order.Floor-currentFloor))) //adds cost for distance from floor
		
		if dir == 1 && order.ButtonType == 1 { //elevator going up order going down
			cost += 3
		} else if dir == -1 && order.ButtonType == 0 { //elevator going down order going up
			cost += 3
		}

		return cost
	
	}else{
		return 101		
	}
	
}
