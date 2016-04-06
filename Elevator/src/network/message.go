package network

//Kan evt flyttes til config - kanskje bedre

const( //type of network message
	IP = iota
	ElevatorPassingFloor
	ElevatorStopping //trenger vi denne?
	ElevatorRunning //trenger vi denne?
	ElevatorAdded
	ElevatorRemoved
	Ack
	//Master
)

type Message struct{
	Source int
	Id int
	Floor int
	Target int //hvilken heis som er target? trenger ikke target for etasje
}