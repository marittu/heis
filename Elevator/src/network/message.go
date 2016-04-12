package network

import (
	"../elevatorDriver"
)

//Kan evt flyttes til config - kanskje bedre

const ( //type of network message
	Init             = 0
	Ping             = 1
	NewOrder         = 2
	OrderFromMaster  = 3
	Ack              = 4
	NewInternalOrder = 5
	Removed          = 6
)

type Message struct {
	MessageId int
	FromIP    string
	ToIP      string
	Order     elevatorDriver.Order
	Info      elevatorDriver.ElevInfo
}
