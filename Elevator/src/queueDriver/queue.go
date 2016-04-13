package queueDriver

import (
	"../elevatorDriver"
	"../network"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

//Queue holds the orders for elevator running on that program
var Queue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS]int{}

//MasterQueue contains all external orders until they are executed
var MasterQueue = [elevatorDriver.N_FLOORS][elevatorDriver.N_BUTTONS - 1]int{}
var Info elevatorDriver.ElevInfo

var mutex sync.Mutex

func QueueInit() {
	FileRead(elevatorDriver.QUEUE)
	time.Sleep(100 * time.Millisecond)
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
	if order.ButtonType < 2 { //only external orders
		MasterQueue[order.Floor][order.ButtonType] = 1
		elevatorDriver.ElevSetButtonLamp(order.Floor, order.ButtonType, 1) //light turned on for all elevators
	}
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
		}
		if button < 2 { //only external orders
			MasterQueue[floor][button] = 0
		}

		for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
			if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
				elevatorDriver.ConnectedElevs[elev].CostQueue[floor][button] = 0
				elevatorDriver.ElevSetButtonLamp(floor, button, 0)
			}
		}
	}

	FileWrite(elevatorDriver.QUEUE)
}

func openDoor(floor int, selfIP string, chToNetwork chan network.Message) {

	DeleteOrder(floor, selfIP)
	setDir(0, selfIP, chToNetwork)
	const duration = 2 * time.Second

	timer := time.NewTimer(0)
	timer.Stop()
	timer.Reset(duration)
	elevatorDriver.ElevSetDoorOpenLamp(1)
	select {
	case <-timer.C:
		timer.Stop()
		elevatorDriver.ElevSetDoorOpenLamp(0)

	}

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

}

func setCurrentFloor(floor int, selfIP string, chToNetwork chan network.Message) {
	mutex.Lock()

	Info.CurrentFloor = floor
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
			elevatorDriver.ConnectedElevs[elev].Info.CurrentFloor = floor
		}
	}
	mutex.Unlock()
	var temp elevatorDriver.ElevInfo
	temp.CurrentFloor = floor
	var msg network.Message
	msg.FromIP = selfIP
	msg.Info = temp
	msg.MessageId = network.Floor

	chToNetwork <- msg
}

func GetCurrentFloor() int {
	mutex.Lock()
	defer mutex.Unlock()
	return Info.CurrentFloor
}

//Gets current floor for a given elevator
func GetInfoIP(IP string) elevatorDriver.ElevInfo {
	mutex.Lock()
	defer mutex.Unlock()
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if elevatorDriver.ConnectedElevs[elev].IP == IP {
			return elevatorDriver.ConnectedElevs[elev].Info
			break
		}
	}

	return elevatorDriver.ElevInfo{CurrentFloor: 0, Dir: 0}
}

func GetDir() int {
	mutex.Lock()
	defer mutex.Unlock()
	return Info.Dir
}

func setDir(dir int, selfIP string, chToNetwork chan network.Message) {
	mutex.Lock()
	Info.Dir = dir
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
			elevatorDriver.ConnectedElevs[elev].Info.Dir = dir
		}
	}
	mutex.Unlock()
	var temp elevatorDriver.ElevInfo
	temp.Dir = dir
	var msg network.Message
	msg.FromIP = selfIP
	msg.Info = temp
	msg.MessageId = network.Dir

	chToNetwork <- msg
}

func PassingFloor(floor int, selfIP string, chToNetwork chan network.Message) {

	setCurrentFloor(floor, selfIP, chToNetwork)
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
		setDir(0, selfIP, chToNetwork)

	} else {
		if Queue[floor][2] == 1 { //internal order
			stopAtFloor(floor, selfIP, chToNetwork)
		} else if dir == 1 && Queue[floor][0] == 1 { //order up, dir up
			stopAtFloor(floor, selfIP, chToNetwork)
		} else if dir == -1 && Queue[floor][1] == 1 { //order down dir down
			stopAtFloor(floor, selfIP, chToNetwork)
		} else if OrderBelow(floor) == false && dir == -1 && Queue[floor][0] == 1 { //order up going down
			stopAtFloor(floor, selfIP, chToNetwork)
		} else if OrderAbove(floor) == false && dir == 1 && Queue[floor][1] == 1 { //order down going up
			stopAtFloor(floor, selfIP, chToNetwork)

		}
	}

}

func stopAtFloor(floor int, selfIP string, chToNetwork chan network.Message) {
	elevatorDriver.ElevDrive(0)
	time.Sleep(100 * time.Millisecond)
	openDoor(floor, selfIP, chToNetwork)
}

func GetDirection(selfIP string, chToNetwork chan network.Message) {
	currentDir := GetDir()
	currentFloor := GetCurrentFloor()
	if EmptyQueue() == true {
		setDir(0, selfIP, chToNetwork)
	} else {

		switch currentDir {
		case 0:
			for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					if Queue[floor][button] == 1 {
						if floor > currentFloor {
							setDir(1, selfIP, chToNetwork)
							elevatorDriver.ElevDrive(1)
						} else if floor < currentFloor {
							setDir(-1, selfIP, chToNetwork)
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
				setDir(-1, selfIP, chToNetwork)
				elevatorDriver.ElevDrive(-1)
			}
		case -1:
			if OrderBelow(currentFloor) {

				elevatorDriver.ElevDrive(-1)
			} else if OrderAbove(currentFloor) {
				setDir(1, selfIP, chToNetwork)
				elevatorDriver.ElevDrive(1)
			}

		}

	}

}

//slett før levering
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
