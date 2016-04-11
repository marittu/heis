package network
import (
	"../elevatorDriver"
)
//Kan evt flyttes til config - kanskje bedre

const( //type of network message
	Ping = 1
	NewOrder = 2
	//Cost = 3
	//FindTarget = 4
	//MasterDistributesOrder = 5
	OrderFromMaster = 3
	NewInternalOrder = 5
	Ack = 4
	
)

type Message struct{
	MessageId int
	FromIP string
	ToIP string
	Order elevatorDriver.Order
	//Cost int
	
}