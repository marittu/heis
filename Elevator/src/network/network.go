package network

import(
	"time"
	"net"
	"strings"
	//"fmt"
	"strconv"
	
)

type Connection struct{
	IP string
	LastPing time.Time
}

var ConnectedElevs []Connection

type elevManager struct{
	SelfIP string
	Elevators Connection
	Master string  
}

var elev elevManager

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
			added := false
			for elev := 0; elev < len(ConnectedElevs); elev++{
				//fmt.Println("Elevators added: ", ConnectedElevs[elev].IP)
				if received.FromIP == ConnectedElevs[elev].IP{
					added = true
				}
				
			}
			if added == false{
				AppendConn(received.FromIP)

			}

			selectMaster()

			chOut <- received

		case send := <-chIn:
			chUDPSend <- send

			
		}
	}

}

func ElevManagerInit() elevManager {

	
	return elev



}

func AppendConn(IP string){
	var temp Connection
	temp.IP = IP
	temp.LastPing = time.Now()

	ConnectedElevs = append(ConnectedElevs, temp)
}

func selectMaster(){
	var masterIP string
	min := 256
	for i, _ := range ConnectedElevs{

	endIP, _ := strconv.Atoi(strings.Replace(ConnectedElevs[i].IP, "129.241.187.", "", -1))
	
		if endIP < min{
			min = endIP
			masterIP = ConnectedElevs[i].IP
		}
	}
	
	elev.Master = masterIP
	
	//fmt.Println("Master: ", elev.Master)
}




//func RemoveConn()