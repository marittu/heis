package userInterfaceDriver

import (
	"../elevatorDriver"
	"time"
	//"fmt"
)

func NewOrder(chButtonPressed chan elevatorDriver.Button){
	for{
		for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
				pressed := elevatorDriver.ElevGetButtonSignal(button, floor)
				if pressed == 1{
					chButtonPressed <- elevatorDriver.Button{ButtonType: button, Floor: floor} //Multiple presses registered?
				}
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