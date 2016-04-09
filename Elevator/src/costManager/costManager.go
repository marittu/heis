package costManager

import(
	"../elevatorDriver"
	//"../network"
	"../queueDriver"
	//"fmt"
)
/*
func GetTargetElev() string{
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(),"/")[0]
	chUDPSend := make(chan Message, 100)
	chUDPReceive := make(chan Message, 100)
	//target := elevatorDriver.ConnectedElevs[0].IP
	
	ownCost := GetOwnCost
	go network.BroadcastCost(SelfIP, ownCost, chUDPSend)
	go network.UDPListener(chUDPReceive)
}
*/
//Finds the best elevator for the order
func GetOwnCost(order elevatorDriver.Order) int{
	//elevator := network.GetElevManager()
	//min := elevatorDriver.N_FLOORS
	var cost int
	//target := elevatorDriver.ConnectedElevs[0].IP

	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{
		
		//fmt.Println("Target IP ",  elevatorDriver.ConnectedElevs[elev].IP)
		if queueDriver.GetCurrentFloor() == order.Floor && queueDriver.GetDir() == 0{ //elevator idle at ordered floor
			cost = 0
			//target = elevatorDriver.ConnectedElevs[elev].IP
			return cost
		}

		if queueDriver.Queue[order.Floor][order.ButtonType] == 1{ //elevator already has a order at the floor
			cost = 1
			//target = elevatorDriver.ConnectedElevs[elev].IP
			return cost
		}

		//cost = getCostForOrder(order. elevatorDriver.Order)	
	}
	return cost
}

/*

//Determines the cost for a given elevator
func getCostForOrder(){

}*/