package queueDriver

import (
	"../elevatorDriver"
	"../network"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

var Queue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{}
var MasterQueue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{}
var Info elevatorDriver.ElevInfo

func QueueInit() {

	FileRead(elevatorDriver.QUEUE)
	MasterQueue = Queue
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS-1; button++ {
			Queue[floor][button] = 0 //Exteral orders handled by other elevators
			elevatorDriver.ElevSetButtonLamp(floor, button, 0)
		}
	}

}

func AddOrder(order elevatorDriver.Order) {
	fmt.Println("Order at: ", order.Floor)
	Queue[order.Floor][order.ButtonType] = 1
	elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1)
	FileWrite(elevatorDriver.QUEUE)
}

func AddOrderMasterQueue(order elevatorDriver.Order) {
	MasterQueue[order.Floor][order.ButtonType] = 1
	elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1)

}

func EmptyQueue() bool {
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if Queue[floor][button] == 1 {
				return false
			}
		}
	}
	return true
}

func OrderAbove(floor int) bool {
	for floor := floor + 1; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if Queue[floor][button] == 1 {
				return true
			}
		}
	}
	return false
}

func OrderBelow(floor int) bool {
	for floor := floor - 1; floor >= 0; floor-- {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if Queue[floor][button] == 1 {
				return true
			}
		}
	}
	return false
}

func DeleteOrder(floor int, selfIP string) {
	for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
		if Queue[floor][button] == 1 {
			Queue[floor][button] = 0

			MasterQueue[floor][button] = 0
		}

		for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
			if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
				elevatorDriver.ConnectedElevs[elev].OwnQueue[floor][button] = 0
				elevatorDriver.ElevSetButtonLamp(floor, button, 0)
			}
		}

	}

	FileWrite(elevatorDriver.QUEUE)
}

func openDoor(floor int, selfIP string, chToNetwork chan network.Message) {

	DeleteOrder(floor, selfIP)
	elevatorDriver.ElevSetDoorOpenLamp(1)
	time.Sleep(2 * time.Second)
	elevatorDriver.ElevSetDoorOpenLamp(0)
	setDir(0, selfIP)

	var temp elevatorDriver.ElevInfo
	temp.Dir = 0
	temp.CurrentFloor = floor
	var msg network.Message
	msg.ToIP = elevatorDriver.ConnectedElevs[0].Master
	msg.FromIP = selfIP
	msg.Info = temp
	msg.MessageId = network.Ack

	chToNetwork <- msg
	GetDirection(selfIP, chToNetwork)
	//printQueue()

}

func setCurrentFloor(floor int, selfIP string) {
	Info.CurrentFloor = floor
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
			elevatorDriver.ConnectedElevs[elev].Info.CurrentFloor = floor
		}
	}
}

func GetCurrentFloor() int {

	return Info.CurrentFloor
}

func GetCurrentFloorIP(IP string) int {
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if elevatorDriver.ConnectedElevs[elev].IP == IP {
			return elevatorDriver.ConnectedElevs[elev].Info.CurrentFloor
		}
	}

	return -1
}

func GetDir() int {
	return Info.Dir
}

func setDir(dir int, selfIP string) {
	Info.Dir = dir

	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
			elevatorDriver.ConnectedElevs[elev].Info.Dir = dir
		}
	}

}

func PassingFloor(floor int, selfIP string, chToNetwork chan network.Message) {

	/*for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if selfIP == elevatorDriver.ConnectedElevs[elev].IP {
			for Floor := 0; Floor < elevatorDriver.N_FLOORS; Floor++ {
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					fmt.Print(elevatorDriver.ConnectedElevs[elev].OwnQueue[Floor][button])
				}
				fmt.Println()
			}
			fmt.Println()
		}
	}*/
	PrintQueue()
	setCurrentFloor(floor, selfIP)
	elevatorDriver.ElevSetFloorIndicator(floor)
	dir := GetDir()

	if floor == 0 {
		elevatorDriver.ElevDrive(0)
		time.Sleep(100 * time.Millisecond)
		GetDirection(selfIP, chToNetwork)
	} else if floor == 3 {
		elevatorDriver.ElevDrive(0)
		time.Sleep(100 * time.Millisecond)
		GetDirection(selfIP, chToNetwork)
	}

	if EmptyQueue() == true {
		elevatorDriver.ElevDrive(0)
		setDir(0, selfIP)
		var temp elevatorDriver.ElevInfo
		temp.Dir = 0
		temp.CurrentFloor = floor
		var msg network.Message
		msg.Info = temp
		msg.MessageId = network.Ping

		chToNetwork <- msg

	} else {
		if Queue[floor][2] == 1 { //internal order
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP, chToNetwork)

		} else if dir == 1 && Queue[floor][0] == 1 { //order up, dir up
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP, chToNetwork)

		} else if dir == -1 && Queue[floor][1] == 1 { //order down dir down
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP, chToNetwork)

		} else if OrderBelow(floor) == false && dir == -1 && Queue[floor][0] == 1 { //order up going down
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP, chToNetwork)

		} else if OrderAbove(floor) == false && dir == 1 && Queue[floor][1] == 1 { //order down going up
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
			openDoor(floor, selfIP, chToNetwork)

		} /*else if dir == -1 && floor == 0 {
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)

		} else if dir == 1 && floor == 3 {
			elevatorDriver.ElevDrive(0)
			time.Sleep(100 * time.Millisecond)
		}*/
	}

}

func GetDirection(selfIP string, chToNetwork chan network.Message) {

	currentDir := GetDir()
	currentFloor := GetCurrentFloor()
	if EmptyQueue() == true {
		setDir(0, selfIP)

	} else {

		switch currentDir {
		case 0:
			for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					if Queue[floor][button] == 1 {
						if floor > currentFloor {
							setDir(1, selfIP)
							elevatorDriver.ElevDrive(1)
						} else if floor < currentFloor {
							setDir(-1, selfIP)
							elevatorDriver.ElevDrive(-1)
						} else if floor == currentFloor {
							openDoor(floor, selfIP, chToNetwork)
						}

					}
				}
			}
		case 1:
			if OrderAbove(currentFloor) {
				elevatorDriver.ElevDrive(1)
			} else if OrderBelow(currentFloor) {
				setDir(-1, selfIP)
				elevatorDriver.ElevDrive(-1)
			}
		case -1:
			if OrderBelow(currentFloor) {

				elevatorDriver.ElevDrive(-1)
			} else if OrderAbove(currentFloor) {
				setDir(1, selfIP)
				elevatorDriver.ElevDrive(1)
			}

		}

	}

}

func PrintQueue() {
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
			fmt.Print(Queue[floor][button])
		}
		fmt.Println()
	}
	fmt.Println()
}

func FileRead(file string) {
	input, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error: file not read")
		return
	}
	err = json.Unmarshal(input, &Queue)
	if err != nil {
		fmt.Println("Error: could not read json file", err.Error())
	}
}

func FileWrite(file string) {
	output, err := json.Marshal(Queue)
	if err != nil {
		fmt.Println("Error queue not json encoded")
	}
	err = ioutil.WriteFile(file, output, 0666)
	if err != nil {
		fmt.Println("Error: file not written")
	}
}
