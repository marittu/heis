package network

import (
	"../elevatorDriver"
)

//Kan evt flyttes til config - kanskje bedre

const ( //type of network message

	Ping             = 1
	NewOrder         = 2
	OrderFromMaster  = 3
	Ack              = 4
	NewInternalOrder = 5
	Removed          = 6
	MovingTimeOut 	 = 7
	Floor            = 8 
	Dir              = 9

)

type Message struct {
	MessageId int
	FromIP    string
	ToIP      string
	Order     elevatorDriver.Order
	Info      elevatorDriver.ElevInfo
}
