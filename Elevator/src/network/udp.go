package network //Endre pakkenavn til liten n

import (
	"net"
	"strings"
	"encoding/json"
	"fmt"

)

const (
	PORT = "30005"
)

func UDPSender (chSend chan Message){
	broadcastAddr := []string{"129.241.187.255", PORT}
	broadcastUDP, _ := net.ResolveUDPAddr("udp", strings.Join(broadcastAddr, ""))
	broadcastConn, _ := net.DialUDP("udp", nil, broadcastUDP)
	defer broadcastConn.Close()
	for{
		msg, err := json.Marshal(<- chSend)
		if err != nil{
			fmt.Println(err)	
		}else{
			broadcastConn.Write(msg)
		}
	}
}

func UDPListener(chReceive chan Message){
	UDPReceiveAddr, err := net.ResolveUDPAddr("udp", PORT)
	if err != nil{
		fmt.Println(err)
	}
	UDPConn, err := net.ListenUDP("udp", UDPReceiveAddr)
	if err != nil{
		fmt.Println(err)
	}
	defer UDPConn.Close()
	msg := make([]byte, 2048)
	trimmedMsg := make([]byte, 1)
	var receivedMessage Message
	for{
		n, _, _ := UDPConn.ReadFromUDP(msg)
		trimmedMsg = msg[:n]
		err := json.Unmarshal(trimmedMsg, &receivedMessage)
		if err != nil{
			fmt.Println(err)
		}else{
			chReceive <- receivedMessage
			}

	}

}