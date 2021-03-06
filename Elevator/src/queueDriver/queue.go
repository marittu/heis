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

func emptyQueue() bool {
	for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if Queue[floor][button] == 1 {
				return false
			}
		}
	}
	return true
}

func orderAbove(floor int) bool {
	for floor := floor + 1; floor < elevatorDriver.N_FLOORS; floor++ {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if Queue[floor][button] == 1 {
				return true
			}
		}
	}
	return false
}

func orderBelow(floor int) bool {
	for floor := floor - 1; floor >= 0; floor-- {
		for button := 0; button < elevatorDriver.N_BUTTONS; button++ {
			if Queue[floor][button] == 1 {
				return true
			}
		}
	}
	return false
}

func deleteOrder(floor int, selfIP string) {
	for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
		if Queue[floor][button] == 1 {
			Queue[floor][button] = 0
			elevatorDriver.ElevSetButtonLamp(floor, button, 0)
		}

		for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
			if elevatorDriver.ConnectedElevs[elev].IP == selfIP {
				elevatorDriver.ConnectedElevs[elev].CostQueue[floor][button] = 0
			}
		}
	}

	FileWrite(elevatorDriver.QUEUE)
}

func openDoor(floor int, selfIP string, chToNetwork chan<- network.Message, Doortimer *time.Timer) {
	
	deleteOrder(floor, selfIP)
	setDir(0, selfIP, chToNetwork)
	const duration = 2 * time.Second

	Doortimer.Stop()
	Doortimer.Reset(duration)
	elevatorDriver.ElevSetDoorOpenLamp(1)

	elevatorDriver.Info.State = elevatorDriver.DoorOpen 	
	var temp elevatorDriver.ElevInfo
	temp.Dir = 0
	temp.CurrentFloor = floor
	var msg network.Message

	//If program killed while door open, elevator might not be connected yet
	if len(elevatorDriver.ConnectedElevs) > 1{ 
		msg.ToIP = elevatorDriver.ConnectedElevs[0].Master	
	}else{
		msg.ToIP = selfIP
	}

	msg.FromIP = selfIP
	msg.Info = temp
	msg.MessageId = network.Ack

	chToNetwork <- msg
}

func setCurrentFloor(floor int, selfIP string, chToNetwork chan network.Message) {
	mutex.Lock()

	elevatorDriver.Info.CurrentFloor = floor
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
	return elevatorDriver.Info.CurrentFloor
}

//Gets current floor and direction for a given elevator
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
	return elevatorDriver.Info.Dir
}

func setDir(dir int, selfIP string, chToNetwork chan<- network.Message) {
	mutex.Lock()
	d := elevatorDriver.ElevMotorDirection(dir)
	elevatorDriver.ElevDrive(d)
	elevatorDriver.Info.Dir = dir
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

func PassingFloor(floor int, selfIP string, chToNetwork chan network.Message, Doortimer *time.Timer, MovingTimer *time.Timer) {
	switch(elevatorDriver.Info.State){
	case elevatorDriver.Moving:
		setCurrentFloor(floor, selfIP, chToNetwork)
		elevatorDriver.ElevSetFloorIndicator(floor)
		dir := GetDir()

		if emptyQueue() == true { 
			MovingTimer.Stop()
			setDir(0, selfIP, chToNetwork)
			time.Sleep(100 * time.Millisecond)
			elevatorDriver.Info.State = elevatorDriver.Idle
		} else if emptyQueue() == false{
			if Queue[floor][2] == 1 { //internal order
				stopAtFloor(floor, selfIP, chToNetwork, Doortimer, MovingTimer)
			} else if dir == 1 && Queue[floor][0] == 1 { //order up, dir up
				stopAtFloor(floor, selfIP, chToNetwork, Doortimer, MovingTimer)
			} else if dir == -1 && Queue[floor][1] == 1 { //order down dir down
				stopAtFloor(floor, selfIP, chToNetwork, Doortimer, MovingTimer)
			} else if orderBelow(floor) == false && dir == -1 && Queue[floor][0] == 1 { //order up going down
				stopAtFloor(floor, selfIP, chToNetwork, Doortimer, MovingTimer)
			} else if orderAbove(floor) == false && dir == 1 && Queue[floor][1] == 1 { //order down going up
				stopAtFloor(floor, selfIP, chToNetwork, Doortimer, MovingTimer)
			}
			//Should stop on endfloors
		}else if floor == 0 {
			MovingTimer.Stop()
			setDir(0, selfIP, chToNetwork)
			time.Sleep(100 * time.Millisecond)
			GetNextOrder(selfIP, chToNetwork, Doortimer, MovingTimer)
			elevatorDriver.Info.State = elevatorDriver.Idle
		} else if floor == 3 {
			MovingTimer.Stop()
			setDir(0, selfIP, chToNetwork)
			time.Sleep(100 * time.Millisecond)
			GetNextOrder(selfIP, chToNetwork, Doortimer, MovingTimer)
			elevatorDriver.Info.State = elevatorDriver.Idle
		}

	default:
		elevatorDriver.Info.State = elevatorDriver.Idle
		MovingTimer.Stop()
		elevatorDriver.Info.TimedOut = false
		setDir(0, selfIP, chToNetwork)
		time.Sleep(1000 * time.Millisecond)
		GetNextOrder(selfIP, chToNetwork, Doortimer, MovingTimer)
	}
}

func stopAtFloor(floor int, selfIP string, chToNetwork chan network.Message, Doortimer *time.Timer, MovingTimer *time.Timer) {
	MovingTimer.Stop()
	setDir(0, selfIP, chToNetwork)
	time.Sleep(100 * time.Millisecond)
	openDoor(floor, selfIP, chToNetwork, Doortimer)
}

func GetNextOrder(selfIP string, chToNetwork chan<- network.Message, Doortimer *time.Timer, MovingTimer *time.Timer) {
	currentDir := GetDir()
	currentFloor := GetCurrentFloor()
	if emptyQueue() == true {
		setDir(0, selfIP, chToNetwork)
		elevatorDriver.Info.State = elevatorDriver.Idle
	} else {
		elevatorDriver.Info.State = elevatorDriver.Moving
		const duration = 8 * time.Second
		MovingTimer.Stop()
		MovingTimer.Reset(duration)
		elevatorDriver.Info.TimedOut = false
		switch currentDir {
		case 0:
			for floor := 0; floor < elevatorDriver.N_FLOORS; floor++ {
				for button := elevatorDriver.BUTTON_CALL_UP; button < elevatorDriver.N_BUTTONS; button++ {
					if Queue[floor][button] == 1 {
						if floor == currentFloor {
							openDoor(floor, selfIP, chToNetwork, Doortimer)
						}else if floor > currentFloor {
							setDir(1, selfIP, chToNetwork)
						} else if floor < currentFloor {
							setDir(-1, selfIP, chToNetwork)
						} 
					}
				}
			}
		case 1:
			if orderAbove(currentFloor) {
				setDir(1, selfIP, chToNetwork)
			} else if orderBelow(currentFloor) {
				setDir(-1, selfIP, chToNetwork)
			}
		case -1:
			if orderBelow(currentFloor) {
				setDir(-1, selfIP, chToNetwork)
			} else if orderAbove(currentFloor) {
				setDir(1, selfIP, chToNetwork)
			}
		}
	}
}

//Internal Order backup
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
