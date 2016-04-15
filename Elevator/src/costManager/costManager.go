package costManager

import (
	"../elevatorDriver"
	"../queueDriver"
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
		//higher cost for more orders
		for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
			for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
				if elevatorDriver.ConnectedElevs[pos].CostQueue[floor][button] == 1 {
					cost += 2 
				}
			}
		}
		//adds cost for distance from floor
		cost += 4 * int(math.Abs(float64(order.Floor-currentFloor))) 
		
		//elevator going up order going down
		if dir == 1 && order.ButtonType == 1 { 
			cost += 3
		//elevator going down order going up
		} else if dir == -1 && order.ButtonType == 0 { 
			cost += 3
		}

		return cost
	
	}else{
		//Elevator Timed out
		return 101		
	}
	
}
