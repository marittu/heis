package network

import(

	"../elevatorDriver"
	"../costManager"
	"time"
	"net"
	"strings"
	"fmt"
	"strconv"
	
)

var SELFIP string
var conn map[string]bool
var cost = make(map[string]int)
var elev elevatorDriver.ElevManager
//cost = make(map[string]bool)
func broadcastIP(IP string, chSend chan Message){
	for{
		chSend <- Message{FromIP: IP, MessageId: Ping, ToIP: ""}
		time.Sleep(100*time.Millisecond)
		
	}
}

func BroadcastCost(IP string, order elevatorDriver.Order, chSend chan Message){
	
	cost := costManager.GetOwnCost(order)
	for{
		chSend <- Message{FromIP: IP, MessageId: Cost, Cost: cost, ToIP: ""}
		time.Sleep(100*time.Millisecond)
		
	}
}

func NetworkHandler(chIn chan Message, chOut chan Message){ 
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(),"/")[0]
	
	conn = make(map[string]bool)
	
	chUDPSend := make(chan Message, 100)
	chUDPReceive := make(chan Message, 100)
	//ownCost := costManager.GetOwnCost()
	go broadcastIP(SelfIP, chUDPSend)
	go UDPListener(chUDPReceive)
	go UDPSender(chUDPSend)
	
	

	for{
		select{
		case received := <- chUDPReceive:
			if received.MessageId == Ping{
				appendConn(received.FromIP)
				
				for elevs := 0; elevs < len(elevatorDriver.ConnectedElevs); elevs++{

					if received.FromIP ==  elevatorDriver.ConnectedElevs[elevs].IP{
						elevatorDriver.ConnectedElevs[elevs].LastPing = time.Now()
						
					}
						
					stillAlive := elevatorDriver.ConnectedElevs[elevs]
					
					if (time.Since(stillAlive.LastPing) > 1000*time.Millisecond){
						removeConn(elevs)
					
					}
					
				}
			}

			if received.MessageId == NewOrder{
				fmt.Println("Broscasting")
				go BroadcastCost(SelfIP, received.Order, chUDPSend)
			}
			
			if received.MessageId == Cost{
				
				AppendCost(received.FromIP, received.Cost)
				if len(cost) == len(elevatorDriver.ConnectedElevs){
					chOut <- Message{FromIP: "", MessageId: FindTarget, ToIP: elev.Master}
				}

			}
			break

			
				
			chOut <- received

		case send := <-chIn:

			chUDPSend <- send

			
		}
	}

}

func GetElevManager() elevatorDriver.ElevManager {
	
	return elev

}


func appendConn(IP string){

	if _, ok := conn[IP]; ok{
		//IP already added
	}else{
		var temp elevatorDriver.Connection
			temp.IP = IP
			temp.LastPing = time.Now()
			elev.SelfIP = IP
			
			elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs, temp)
			fmt.Println("Connected elevator: ", IP)
			fmt.Println("")
			conn[IP] = true
			selectMaster()
 	}

}

func AppendCost(IP string, ownCost int){
	if _, ok := cost[IP]; ok{
		//cost already addded
	}else{
		//cost = append(IP, cost)
		cost[IP] = ownCost
		fmt.Println("Cost added: ", IP, " cost: ", ownCost)
			
 	}
 	

}

func GetMinCost() string{
	min := 1000
	var ideal string
	for IP, minCost := range cost{
		if minCost < min{
			ideal = IP
			min = minCost
		}
	}
	//fmt.Println("Ideal: ", ideal)
	return ideal

}

func selectMaster(){
	var masterIP string
	min := 256
	for i, _ := range elevatorDriver.ConnectedElevs{

	endIP, _ := strconv.Atoi(strings.Replace(elevatorDriver.ConnectedElevs[i].IP, "129.241.187.", "", -1))
	
		if endIP < min{
			min = endIP
			masterIP = elevatorDriver.ConnectedElevs[i].IP
		}
	}
	
	elev.Master = masterIP
	
	fmt.Println("Master: ", elev.Master)
}



func removeConn(elev int){
	fmt.Println("Removed: ", elevatorDriver.ConnectedElevs[elev].IP)
	delete(conn, elevatorDriver.ConnectedElevs[elev].IP)
	elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs[:elev], elevatorDriver.ConnectedElevs[elev+1:]...)
	selectMaster()
}

func SendNetworkMessage(order elevatorDriver.Order, selfIP string, toIP string, msgId int, chToNetwork chan Message){
	var msg Message
	msg.Order = order
	msg.ToIP = toIP
	msg.FromIP = selfIP
	msg.MessageId = msgId

	chToNetwork <- msg
}
