package network
import (
	"../elevatorDriver"
)
//Kan evt flyttes til config - kanskje bedre

const( //type of network message
	Ping = 1
	NewOrder = 2
	MasterDistributesOrder = 3
	
)

type Message struct{
	MessageId int
	FromIP string
	ToIP string
	Order elevatorDriver.Order
	
}