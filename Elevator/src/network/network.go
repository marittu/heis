package network

import(

	"../elevatorDriver"
	"time"
	"net"
	"strings"
	"fmt"
	"strconv"
	
)


var conn map[string]bool
var elev elevatorDriver.ElevManager

func broadcastIP(IP string, chSend chan Message){
	for{
		chSend <- Message{FromIP: IP, MessageId: Ping, ToIP: ""}
		time.Sleep(100*time.Millisecond)
	}
}



func NetworkHandler(chIn chan Message, chOut chan Message){ 
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(),"/")[0]
	//fmt.Println(selfIP)
	chUDPSend := make(chan Message, 100)
	chUDPReceive := make(chan Message, 100)
	//var test map[string]IP
	conn = make(map[string]bool)
	go broadcastIP(SelfIP, chUDPSend)
	go UDPListener(chUDPReceive)
	go UDPSender(chUDPSend)
	

	for{
		//fmt.Println("Start of for loop")
		select{
		case received := <- chUDPReceive:
			
			AppendConn(received.FromIP)
			selectMaster()			
			for elevs := 0; elevs < len(elevatorDriver.ConnectedElevs); elevs++{

				fmt.Println("Connected elevators: ", elevatorDriver.ConnectedElevs[elevs].IP)
				//selectMaster()		
				fmt.Println("Master: ", elev.Master)
				

				if received.MessageId == Ping{
					if received.FromIP ==  elevatorDriver.ConnectedElevs[elevs].IP{
						elevatorDriver.ConnectedElevs[elevs].LastPing = time.Now()	
					}
					
					
				}
				

				stillAlive := elevatorDriver.ConnectedElevs[elevs]
				
				if (time.Since(stillAlive.LastPing) > 900*time.Millisecond){
					fmt.Println("Removing Connection: ", stillAlive.IP)
					RemoveConn(elevs)
				
				}
				
				

			}
			chOut <- received

		case send := <-chIn:
			fmt.Println("Sending over chUDPSend")
			chUDPSend <- send
			fmt.Println("Done chUDPSend")

			
		}
	}

}

func ElevManagerInit() elevatorDriver.ElevManager {

	selectMaster()
	
	return elev



}


func AppendConn(IP string){

	if _, ok := conn[IP]; ok{
		
	}else{
		var temp elevatorDriver.Connection
			temp.IP = IP
			temp.LastPing = time.Now()

			elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs, temp)
			fmt.Println("Connected elevators: ", IP)

			conn[IP] = true
 	}

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
	
	//fmt.Println("Master: ", elev.Master)
}



func RemoveConn(elev int){
	fmt.Println("Removed: ", elevatorDriver.ConnectedElevs[elev].IP)
	delete(conn, elevatorDriver.ConnectedElevs[elev].IP)
	elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs[:elev], elevatorDriver.ConnectedElevs[elev+1:]...)
	selectMaster()
}