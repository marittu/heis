package elevatorDriver

import (
	"time"
)
 
var eInfo ElevInfo

const (
	PORT = ":30106"
)

const N_FLOORS = 4
const N_BUTTONS = 3

const QUEUE = "queue.txt"

type ElevButtonType int

const (
	BUTTON_CALL_UP   ElevButtonType = 0
	BUTTON_CALL_DOWN                = 1
	BUTTON_INTERNAL                 = 2
)

type ElevMotorDirection int

const (
	DIR_DOWN 					= -1
	DIR_UP                      = 1
	DIR_STOP                    = 0
)

type Order struct {
	ButtonType ElevButtonType
	Floor      int
}

type ElevInfo struct {
	Dir          int
	CurrentFloor int
	State        StateType
	TimedOut 	 bool
}

var Info ElevInfo

type StateType int
const(
	Idle 		= 0
	Moving 		= 1
	DoorOpen 	= 2
)

type Connection struct {
	IP        string
	LastPing  time.Time
	Master    string
	CostQueue [N_FLOORS][N_BUTTONS]int
	Info      ElevInfo
}

var ConnectedElevs []Connection
