package elevatorDriver

import(
	"time"
)


var eInfo ElevInfo


const (
	PORT = ":30105"
)


const N_FLOORS  	= 	4 
const N_BUTTONS		= 	3   

type ElevButtonType int
const (
	BUTTON_CALL_UP ElevButtonType = 0
	BUTTON_CALL_DOWN = 1
	BUTTON_INTERNAL = 2
)

type ElevMotorDirection int
const (
	DIR_DOWN ElevMotorDirection = -1
	DIR_UP = 1
	DIR_STOP = 0
)

type Order struct{  
	ButtonType ElevButtonType
	Floor int
}

type ElevStatus struct {
	Dir ElevMotorDirection
	LastFloor int
}

type ElevInfo struct{
	Dir int
	CurrentFloor int
}


type Connection struct{
	IP string
	LastPing time.Time
}

var ConnectedElevs []Connection

type ElevManager struct{
	SelfIP string
	Elevators Connection
	Master string  
}