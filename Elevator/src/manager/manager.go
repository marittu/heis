package manager

import (
	
	"../elevatorDriver"
	"../queueDriver"
	//"../Network"
	"fmt"
	"sync"
	"time"
)

var Info elevatorDriver.ElevInfo

var mutex sync.Mutex

func setCurrentFloor(floor int){
	mutex.Lock()
	Info.CurrentFloor = floor
    defer mutex.Unlock()
}

func getCurrentFloor() int{
	mutex.Lock()
    defer mutex.Unlock()

	return Info.CurrentFloor
}

func getDir()int{
	return Info.Dir
}

func setDir(dir int){
	Info.Dir = dir
	
}

func ChannelHandler(chButtonPressed chan elevatorDriver.Button, chGetFloor chan int){
	for{ 
		select{
		case order := <- chButtonPressed: //add case for internal order or external order
			queueDriver.AddOrder(order)
			GetDirection()
			break
		case floor := <- chGetFloor:
			setCurrentFloor(floor)
			PassingFloor(floor)
			break
		}
	}
}




func openDoor(floor int){
	queueDriver.DeleteOrder(floor)
	elevatorDriver.ElevSetDoorOpenLamp(1)
	time.Sleep(2*time.Second)
	elevatorDriver.ElevSetDoorOpenLamp(0)
	GetDirection()
	//printQueue()


}

func PassingFloor(floor int){ 
	elevatorDriver.ElevSetFloorIndicator(floor)
	dir := getDir()

	if queueDriver.EmptyQueue() == true{
		elevatorDriver.ElevDrive(0)
		setDir(0)
		
	}else{
		if queueDriver.Queue[floor][2] == 1{
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor)	

		}else if (dir == 1 && queueDriver.Queue[floor][0] == 1){
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor)
			
		}else if (dir == -1 && queueDriver.Queue[floor][1] == 1){
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor)
			
		}	
	}
	
}

func GetDirection(){
	
	currentDir := getDir()
	currentFloor := getCurrentFloor()		
	if queueDriver.EmptyQueue() == true{
		setDir(0)
		
	}else{
		
		switch(currentDir){
		case 0:
			for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					if queueDriver.Queue[floor][button] == 1{
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
			if queueDriver.OrderAbove(currentFloor){
				elevatorDriver.ElevDrive(1)
			}else if queueDriver.OrderBelow(currentFloor){	
				setDir(-1)
				elevatorDriver.ElevDrive(-1)
			}
		case -1:
			if queueDriver.OrderBelow(currentFloor){

				elevatorDriver.ElevDrive(-1)
			}else if queueDriver.OrderAbove(currentFloor){	
				setDir(1)
				elevatorDriver.ElevDrive(1)
			}




		}

	}


}

func printQueue(){
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++{
			for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
				fmt.Print(queueDriver.Queue[floor][button])
			} 
			fmt.Println()
	}
}
