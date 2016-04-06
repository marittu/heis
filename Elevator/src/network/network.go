package network

import(
	"time"
	"net"
	"strings"
	"fmt"
	
)

type Connection struct{
	IP string
	LastPing time.Time
}

var ConnectedElevs []Connection

type elevManager struct{
	selfIP string
	elevators Connection
	master string  
}

var elev elevManager

func broadcastIP(IP string, chSend chan Message){
	for{
		chSend <- Message{FromIP: IP, MessageId: Ping, ToIP: ""}
		time.Sleep(100*time.Millisecond)
	}
}

func Manager(chIn chan Message, chOut chan Message){ 
	addr, _ := net.InterfaceAddrs()
	selfIP := strings.Split(addr[1].String(),"/")[0]
	//fmt.Println(selfIP)
	chUDPSend := make(chan Message, 100)
	chUDPReceive := make(chan Message, 100)

	go broadcastIP(selfIP, chUDPSend)
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
			
		}
	}

}


func AppendConn(IP string){
	var temp Connection
	temp.IP = IP
	temp.LastPing = time.Now()

	ConnectedElevs = append(ConnectedElevs, temp)
}

func selectMaster(){
	min := 256
	for i, _ := range ConnectedElevs{
		if i < min{
			min = i
		}
	}
	elev.master = ConnectedElevs[min].IP
	fmt.Println("Master: ", elev.master)
}




//func RemoveConn()