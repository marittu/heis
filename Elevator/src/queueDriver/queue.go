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

func AddOrder(order elevatorDriver.Order, selfIP string){
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
		if elevatorDriver.ConnectedElevs[elev].IP == selfIP{
			elevatorDriver.ConnectedElevs[elev].OwnQueue[order.Floor][order.ButtonType] = 1
			elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1)
			/*for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
					for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
						fmt.Print(elevatorDriver.ConnectedElevs[elev].OwnQueue[floor][button])
					} 
					fmt.Println()
			}
			fmt.Println()
			*/	
		}
	}
	Queue[order.Floor][order.ButtonType] = 1
	
	

}

func AddOrderMasterQueue(order elevatorDriver.Order){
	MasterQueue[order.Floor][order.ButtonType] = 1
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

func DeleteOrder(floor int, selfIP string){
	for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++{
		
		Queue[floor][button] = 0
		MasterQueue[floor][button] = 0
		for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
			if elevatorDriver.ConnectedElevs[elev].IP == selfIP{
				elevatorDriver.ConnectedElevs[elev].OwnQueue[floor][button] = 0
			}
		}
		
		elevatorDriver.ElevSetButtonLamp(floor,button,0) //send lights over network
	}
}

func openDoor(floor int, selfIP string){
	DeleteOrder(floor, selfIP)
	elevatorDriver.ElevSetDoorOpenLamp(1)
	time.Sleep(2*time.Second)
	elevatorDriver.ElevSetDoorOpenLamp(0)
	GetDirection(selfIP)
	//printQueue()


}

func setCurrentFloor(floor int, selfIP string){
	Info.CurrentFloor = floor
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
			if elevatorDriver.ConnectedElevs[elev].IP == selfIP{
				elevatorDriver.ConnectedElevs[elev].Info.CurrentFloor = floor
			}
	}  
}

func GetCurrentFloor() int{

	return Info.CurrentFloor
}

func GetDir()int{
	return Info.Dir
}

func setDir(dir int, selfIP string){
	Info.Dir = dir

	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
			if elevatorDriver.ConnectedElevs[elev].IP == selfIP{
				elevatorDriver.ConnectedElevs[elev].Info.Dir = dir
			}
	} 
	
}

func PassingFloor(floor int, selfIP string){ 
	setCurrentFloor(floor, selfIP)
	elevatorDriver.ElevSetFloorIndicator(floor)
	dir := GetDir()

	if EmptyQueue() == true{
		elevatorDriver.ElevDrive(0)
		setDir(0, selfIP)
		
	}else{
		if Queue[floor][2] == 1{
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP)	

		}else if (dir == 1 && Queue[floor][0] == 1){
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP)
			
		}else if (dir == -1 && Queue[floor][1] == 1){
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP)
			
		}	
	}
	
}

func GetDirection(selfIP string){
	
	currentDir := GetDir()
	currentFloor := GetCurrentFloor()		
	if EmptyQueue() == true{
		setDir(0, selfIP)
		
	}else{
		
		switch(currentDir){
		case 0:
			for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					if Queue[floor][button] == 1{
						if floor > currentFloor{
							setDir(1, selfIP)
							elevatorDriver.ElevDrive(1)
						}else if floor < currentFloor{
							setDir(-1, selfIP)
							elevatorDriver.ElevDrive(-1)
						}else if floor == currentFloor{
							openDoor(floor, selfIP)
						}

					}
				}
			}
		case 1:
			if OrderAbove(currentFloor){
				elevatorDriver.ElevDrive(1)
			}else if OrderBelow(currentFloor){	
				setDir(-1, selfIP)
				elevatorDriver.ElevDrive(-1)
			}
		case -1:
			if OrderBelow(currentFloor){

				elevatorDriver.ElevDrive(-1)
			}else if OrderAbove(currentFloor){	
				setDir(1, selfIP)
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


