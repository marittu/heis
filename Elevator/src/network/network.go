package network

import(
	"time"
	"net"
	"strings"
	"fmt"
	
)

func broadcastIP(IP string, chSend chan Message){
	for{
		chSend <- Message{FromIP: IP, MessageId: Ping, ToIP: ""}
		time.Sleep(100*time.Millisecond)
	}
}

func Manager(chIn chan Message, chOut chan Message){ //chIn chan Message, chOut chan Message
	addr, _ := net.InterfaceAddrs()
	selfIP := strings.Split(addr[1].String(),"/")[0]
	fmt.Println(selfIP)
	chUDPSend := make(chan Message, 100)
	chUDPReceive := make(chan Message, 100)

	go broadcastIP(selfIP, chUDPSend)
	go UDPListener(chUDPReceive)
	go UDPSender(chUDPSend)

	for{
		select{
		case received := <- chUDPReceive:
			fmt.Println(received.FromIP)
		}
	}
	/*

	for{
		select{
		case msg := <- chUDPReceive:
			_, present := connected[msg.FromIP]

			if msg.MessageId == Ping{ 
				if msg.FromIP != selfIP{
					if present{
						connected[msg.FromIP].LastSignal = time.Now()
					}else{
						if time.Since(connected[msg.FromIP].LastSignal) > 900 * time.Millisecond{
							removeElev(msg.FromIP, chIn)
							chOut <- Message{FromIP : msg.FromIP, MessageId : ElevatorAdded}	
						}
						
					}
				}
				break
			}

			chOut <- msg

		case msg := <-chIn:
			chUDPSend <- msg
		}
	}*/
}
/*
func removeElev(IP string, chOut chan Message){
	//delete(elevTimers, IP)
	chOut <- Message{FromIP: IP, MessageId: ElevatorRemoved}
}*/
