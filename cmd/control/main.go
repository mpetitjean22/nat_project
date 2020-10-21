package main

import (
	"encoding/binary"
	"fmt"
	"nat_project/pkg/control_packet"
	"net"
	"os"
	"strconv"
)

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) != 5 {
		fmt.Println("Invalid Number of Arguments")
		fmt.Println("Looking for: (sourceIP) (sourcePort) (destinationIP) (destinationPort)")
	}

	srcIP := net.ParseIP(argsWithProg[1])
	if srcIP == nil {
		fmt.Printf("Invalid Source IP \n")
		return
	}

	srcPortVal, err := strconv.ParseUint(argsWithProg[2], 10, 16)
	if err != nil {
		fmt.Printf("Invalid Source Port %v \n", err)
		return
	}

	dstIP := net.ParseIP(argsWithProg[3])
	if dstIP == nil {
		fmt.Printf("Invalid Dest IP \n")
		return
	}

	dstPortVal, err := strconv.ParseUint(argsWithProg[4], 10, 16)
	if err != nil {
		fmt.Printf("Invalid Dest Port %v \n", err)
		return
	}

	srcPort := make([]byte, 2)
	binary.BigEndian.PutUint16(srcPort, uint16(srcPortVal))

	dstPort := make([]byte, 2)
	binary.BigEndian.PutUint16(dstPort, uint16(dstPortVal))

	//fmt.Println(srcIP, srcPort, dstIP, dstPort)
	//srcIP := []byte{0x01, 0x01, 0x01, 0x01} // 1.1.1.1
	//dstIP := []byte{0x02, 0x02, 0x02, 0x02} // 2.2.2.2
	//srcPort := []byte{0x00, 0x50} // 80
	//dstPort := []byte{0x00, 0x50} // 80

	// add a mapping from 1.1.1.1:80 to 2.2.2.2:80
	control_packet.SendAddMapping(srcIP[12:16], dstIP[12:16], srcPort, dstPort)

	// list the mapping (this will just print it out)
	control_packet.SendListMappings()
}
