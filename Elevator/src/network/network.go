package network

import (
	"../elevatorDriver"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var conn map[string]bool

func broadcastIP(IP string, chSend chan Message) {
	for {
		chSend <- Message{FromIP: IP, MessageId: Ping, ToIP: ""}
		time.Sleep(100 * time.Millisecond)

	}
}

func NetworkHandler(chIn chan Message, chOut chan Message) {
	addr, _ := net.InterfaceAddrs()
	SelfIP := strings.Split(addr[1].String(), "/")[0]

	conn = make(map[string]bool)

	chUDPSend := make(chan Message, 100)
	chUDPReceive := make(chan Message, 100)
	go broadcastIP(SelfIP, chUDPSend)
	go UDPListener(chUDPReceive)
	go UDPSender(chUDPSend)

	for {
		select {
		case received := <-chUDPReceive:

			if received.MessageId == Init {

				for elevs := 0; elevs < len(elevatorDriver.ConnectedElevs); elevs++ {

					if received.FromIP == elevatorDriver.ConnectedElevs[elevs].IP {
						elevatorDriver.ConnectedElevs[elevs].Info.CurrentFloor = received.Info.CurrentFloor

					}
				}
			}

			if received.MessageId == Ping {
				appendConn(received.FromIP)

				for elevs := 0; elevs < len(elevatorDriver.ConnectedElevs); elevs++ {

					if received.FromIP == elevatorDriver.ConnectedElevs[elevs].IP {
						elevatorDriver.ConnectedElevs[elevs].LastPing = time.Now()

					}
					//fmt.Println(elevatorDriver.ConnectedElevs[elevs].IP, "	", elevatorDriver.ConnectedElevs[elevs].LastPing)
					fmt.Println(elevatorDriver.ConnectedElevs[elevs].IP, "	", elevatorDriver.ConnectedElevs[elevs].Info.CurrentFloor)

					stillAlive := elevatorDriver.ConnectedElevs[elevs]

					if time.Since(stillAlive.LastPing) > 600*time.Millisecond {
						removeConn(elevs, chOut)

					}

				}
			}

			if received.MessageId == NewInternalOrder || received.MessageId == OrderFromMaster {

				for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
					if received.FromIP == elevatorDriver.ConnectedElevs[elev].IP {
						elevatorDriver.ConnectedElevs[elev].OwnQueue[received.Order.Floor][received.Order.ButtonType] = 1
						elevatorDriver.ConnectedElevs[elev].Info = received.Info
						chOut <- received
					}

				}

			}

			if received.MessageId == Ack {

				for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
					if received.FromIP == elevatorDriver.ConnectedElevs[elev].IP {
						elevatorDriver.ConnectedElevs[elev].OwnQueue[received.Order.Floor][received.Order.ButtonType] = 0
						chOut <- received
					}

				}

			}

			chOut <- received

		case send := <-chIn:
			chUDPSend <- send

		}
	}

}

func appendConn(IP string) {

	if _, ok := conn[IP]; ok {
		//IP already added
	} else {
		var temp elevatorDriver.Connection
		temp.IP = IP
		temp.LastPing = time.Now()

		elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs, temp)
		fmt.Println("Connected elevator: ", IP)

		conn[IP] = true
		selectMaster()

	}

}

func selectMaster() {
	var masterIP string
	min := 256
	for i, _ := range elevatorDriver.ConnectedElevs {

		endIP, _ := strconv.Atoi(strings.Replace(elevatorDriver.ConnectedElevs[i].IP, "129.241.187.", "", -1))

		if endIP < min {
			min = endIP
			masterIP = elevatorDriver.ConnectedElevs[i].IP
		}
	}
	for elev := 0; elev < len(elevatorDriver.ConnectedElevs); elev++ {
		elevatorDriver.ConnectedElevs[elev].Master = masterIP
	}

	fmt.Println("Master: ", masterIP)
}

func removeConn(elev int, chOut chan Message) {
	fmt.Println("Removed: ", elevatorDriver.ConnectedElevs[elev].IP)
	delete(conn, elevatorDriver.ConnectedElevs[elev].IP)
	elevatorDriver.ConnectedElevs = append(elevatorDriver.ConnectedElevs[:elev], elevatorDriver.ConnectedElevs[elev+1:]...)
	selectMaster()

	var temp Message
	temp.MessageId = Removed

	chOut <- temp
}
