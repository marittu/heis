package queueDriver

import (
	"../elevatorDriver"
	//"fmt"
)

var Queue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{}

func QueueInit(){
	Queue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{
		{0, -1, 0}, 
		{0, 0, 0}, 
		{0, 0, 0}, 
		{-1, 0, 0}}
}

func AddOrder(order elevatorDriver.Button){
	Queue[order.Floor][order.ButtonType] = 1
	elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1)
}
func EmptyQueue()bool{
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
					if Queue[floor][button] == 1{
						return false
					}
			}
	}
	return true
}

func OrderAbove(floor int)bool{
	for floor := floor + 1; floor < elevatorDriver.N_FLOORS; floor++{
		for button := 0; button < elevatorDriver.N_BUTTONS; button++{
			if Queue[floor][button] == 1{
				return true
			}
		}
	}
	return false
}

func OrderBelow(floor int)bool{
	for floor := floor - 1; floor >= 0; floor--{
		for button := 0; button < elevatorDriver.N_BUTTONS; button++{
			if Queue[floor][button] == 1{
				return true
			}
		}
	}
	return false
}

func DeleteOrder(floor int){
	for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++{
			if Queue[floor][button] == 1{
				Queue[floor][button] = 0
				elevatorDriver.ElevSetButtonLamp(floor,button,0)
			}
		}
}

