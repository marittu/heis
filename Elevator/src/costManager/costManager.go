package costManager

import (
	"../elevatorDriver"
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

		fmt.Println("Cost: ", cost, "for elev: ", elevatorDriver.ConnectedElevs[elev].IP)
	}
	return target
}

func getOwnCost(pos int, order elevatorDriver.Order) int {

	cost := 0

	//elevator already at floor

	if (elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor == order.Floor) && (elevatorDriver.ConnectedElevs[pos].Info.Dir == 0) {
		cost = 0
		return cost
	}

	//elevator already has orders at floor
	for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
		if elevatorDriver.ConnectedElevs[pos].OwnQueue[order.Floor][button] == 1 {
			cost = 1
			return cost
		}
	}

	//higher cost for more orders
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if elevatorDriver.ConnectedElevs[pos].OwnQueue[floor][button] == 1 {
				cost += 1
				fmt.Println("Already has orders")
			}
		}
	}

	cost += int(math.Abs(float64(order.Floor - elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor)))
	fmt.Println("Cost from different floor: ", int(math.Abs(float64(order.Floor-elevatorDriver.ConnectedElevs[pos].Info.CurrentFloor))))

	if elevatorDriver.ConnectedElevs[pos].Info.Dir == 0 { //elevator at floor
		cost += 1
		fmt.Println("Cost from dir")
	} else if elevatorDriver.ConnectedElevs[pos].Info.Dir == 1 && order.ButtonType == 1 { //elevator going up order going down
		cost += 3
	} else if elevatorDriver.ConnectedElevs[pos].Info.Dir == -1 && order.ButtonType == 0 { //elevator going down order going up
		cost += 3
	}

	return cost

}
