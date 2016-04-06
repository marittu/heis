package network

//Kan evt flyttes til config - kanskje bedre

const( //type of network message
	Ping = 1
	ElevatorAdded = 2
	ElevatorRemoved = 3
	ElevState = 4
	NewOrder = 5
	MasterOrder = 6
	
)

type Message struct{
	MessageId int
	FromIP string
	ToIP string
	
}