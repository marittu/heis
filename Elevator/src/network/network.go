package network

import(
	"time"
	"net"
	"fmt"
)

var elevTimers map[int]*time.Timer 

func broadcastIP (ID int, chSend chan Message){
	for{
		chSend <- Message{Source: ID, Id: IP}
		time.Sleep(100*time.Millisecond)
	}
}

func manager(chIn chan Messagen chOut chan Message){
	addr, _ := net.InterfaceAddrs()
	selfID := int(addr[1].Strings()[12] - '0') * 100 + int(addr[1].Strings()[][13] - '0') * 10 + int(addr[1].String()[14] - '0')

	chUDPSend := nake(chan Message, 100)
	chUDPReceive := make(chan Message, 100)

	go broadcastIP(selfID, chUDPSend)
	go UDPLIstener(chUDPReceive)
	go UDPSender(chUDPSend)

	elevTimers = make(map[int]*time.Timer)

	for{
		select{
		case msg := chUDPReceive:
			_, present := elevTimers[msg.Source]

			if msg.ID == IP{ //hvor har vi IP fra
				if msg.Source != selfID{
					if present{
						elevTimers[msg.Source].Reset(time.Second)
					}else{
						elevTimers[msg.Source] = time.AfterFunc(time.Second, func() {removeElev(msg.Source, chIn)})
						chOut <- Message{Source : msg.Source, Id : ElevatorAdded}
					}
				}
				break
			}

			chOut <- msg

		case msg := <-chIn:
			chUDPSend <- msg
		}
	}
}

func removeElev(id int, chOut chan Message){
	delete(elevTimers, id)
	chOut <- Message{Source: id, Id: ElevatorRemoved}
}