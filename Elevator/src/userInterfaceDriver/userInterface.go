package userInterfaceDriver

import (
	"../elevatorDriver"
	//"../queueDriver"
	"time"
	//"fmt"
)

func NewOrder(chButtonPressed chan elevatorDriver.Order){
	var temp [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int
	for{
		time.Sleep(100 * time.Millisecond)
		for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
				pressed := elevatorDriver.ElevGetButtonSignal(button, floor)
				if pressed != 0 && pressed != temp[floor][button]{
					chButtonPressed <- elevatorDriver.Order{ButtonType: button, Floor: floor} 
				} 
				temp[floor][button] = pressed
			}
		}	
	}
	
}

func FloorTracker(chGetFloor chan int){
	var currentFloor int
	previousFloor := -1
	for{
		time.Sleep(100 * time.Millisecond)
		currentFloor = elevatorDriver.ElevGetFloorSensorSignal()
		if currentFloor != previousFloor && currentFloor != -1{
			previousFloor = currentFloor
			chGetFloor <- currentFloor 
		}
	}

}