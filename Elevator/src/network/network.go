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
						
			for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++{

				fmt.Println("Elevators added: ", elevatorDriver.ConnectedElevs[elev].IP)
				selectMaster()

				if received.MessageId == Ping{ 
					elevatorDriver.ConnectedElevs[elev].LastPing = time.Now()
				}
				

				stillAlive := elevatorDriver.ConnectedElevs[elev]
				if (time.Since(stillAlive.LastPing) > 900*time.Millisecond){
					RemoveConn(elev)
				
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
	
	fmt.Println("Master: ", elev.Master)
}



func RemoveConn(elev int){
	elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs[:elev], elevatorDriver.ConnectedElevs[elev+1:]...)
	selectMaster()
}