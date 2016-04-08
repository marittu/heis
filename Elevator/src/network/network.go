package network

import(

	"../elevatorDriver"
	"time"
	"net"
	"strings"
	"fmt"
	"strconv"
	
)



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

	go broadcastIP(SelfIP, chUDPSend)
	go UDPListener(chUDPReceive)
	go UDPSender(chUDPSend)
	

	for{
		select{
		case received := <- chUDPReceive:
			
			AppendConn(received.FromIP)
			selectMaster()			
			for elevs := 0; elevs < len(elevatorDriver.ConnectedElevs); elevs++{

				fmt.Println("Connected elevators: ", elevatorDriver.ConnectedElevs[elevs].IP)
				fmt.Println("Master: ", elev.Master)
				//selectMaster()

				if received.MessageId == Ping{ 
					elevatorDriver.ConnectedElevs[elevs].LastPing = time.Now()
				}
				

				stillAlive := elevatorDriver.ConnectedElevs[elevs]
				if (time.Since(stillAlive.LastPing) > 900*time.Millisecond){
					fmt.Println("Removing Connection: ", stillAlive.IP)
					RemoveConn(elevs)
				
				}
				
				

			}
			chOut <- received

		case send := <-chIn:
			chUDPSend <- send

			
		}
	}

}

func ElevManagerInit() elevatorDriver.ElevManager {

	selectMaster()
	
	return elev



}

func AppendConn(IP string){
	var temp elevatorDriver.Connection
	temp.IP = IP
	temp.LastPing = time.Now()

	elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs, temp)
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
	elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs[:elev], elevatorDriver.ConnectedElevs[elev+1:]...)

	selectMaster()
}