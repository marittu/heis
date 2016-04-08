package queueDriver

import (
	"../elevatorDriver"
	"fmt"
	"time"
)

var Queue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{}
var MasterQueue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{}
var Info elevatorDriver.ElevInfo


func QueueInit(){
	Queue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{
		{0, -1, 0}, 
		{0, 0, 0}, 
		{0, 0, 0}, 
		{-1, 0, 0}}

	MasterQueue = Queue
}

func AddOrder(order elevatorDriver.Button){
	Queue[order.Floor][order.ButtonType] = 1
	elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1)
}

func AddOrderMasterQueue(order elevatorDriver.Button){
	MasterQueue[order.Floor][order.ButtonType] = 1
	elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1)
	//fmt.Println("Fucking up in AddOrderMasterQueue")
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

func openDoor(floor int){
	DeleteOrder(floor)
	elevatorDriver.ElevSetDoorOpenLamp(1)
	time.Sleep(2*time.Second)
	elevatorDriver.ElevSetDoorOpenLamp(0)
	GetDirection()
	//printQueue()


}

func setCurrentFloor(floor int){
	Info.CurrentFloor = floor    
}

func getCurrentFloor() int{
	return Info.CurrentFloor
}

func getDir()int{
	return Info.Dir
}

func setDir(dir int){
	Info.Dir = dir
	
}

func PassingFloor(floor int){ 
	setCurrentFloor(floor)
	elevatorDriver.ElevSetFloorIndicator(floor)
	dir := getDir()

	if EmptyQueue() == true{
		elevatorDriver.ElevDrive(0)
		setDir(0)
		
	}else{
		if Queue[floor][2] == 1{
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor)	

		}else if (dir == 1 && Queue[floor][0] == 1){
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor)
			
		}else if (dir == -1 && Queue[floor][1] == 1){
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor)
			
		}	
	}
	
}

func GetDirection(){
	
	currentDir := getDir()
	currentFloor := getCurrentFloor()		
	if EmptyQueue() == true{
		setDir(0)
		
	}else{
		
		switch(currentDir){
		case 0:
			for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					if Queue[floor][button] == 1{
						if floor > currentFloor{
							setDir(1)
							elevatorDriver.ElevDrive(1)
						}else if floor < currentFloor{
							setDir(-1)
							elevatorDriver.ElevDrive(-1)
						}else if floor == currentFloor{
							openDoor(floor)
						}

					}
				}
			}
		case 1:
			if OrderAbove(currentFloor){
				elevatorDriver.ElevDrive(1)
			}else if OrderBelow(currentFloor){	
				setDir(-1)
				elevatorDriver.ElevDrive(-1)
			}
		case -1:
			if OrderBelow(currentFloor){

				elevatorDriver.ElevDrive(-1)
			}else if OrderAbove(currentFloor){	
				setDir(1)
				elevatorDriver.ElevDrive(1)
			}




		}

	}


}

func PrintQueue(){
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
				fmt.Print(MasterQueue[floor][button])
			} 
			fmt.Println()
	}
	fmt.Println()
}


