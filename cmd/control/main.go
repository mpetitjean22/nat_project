package main

import (
	"encoding/binary"
	"fmt"
	"nat_project/pkg/control_packet"
	"net"
	"os"
	"strconv"
)

func addMapping(argsWithProg []string, controlType int) {
	if len(argsWithProg) != 6 {
		fmt.Println("Invalid Number of Arguments")
		fmt.Println("Looking for: (control type) (sourceIP) (sourcePort) (destinationIP) (destinationPort)")
	}

	srcIP := net.ParseIP(argsWithProg[2])
	if srcIP == nil {
		fmt.Printf("Invalid Source IP \n")
		return
	}

	srcPortVal, err := strconv.ParseUint(argsWithProg[3], 10, 16)
	if err != nil {
		fmt.Printf("Invalid Source Port %v \n", err)
		return
	}

	dstIP := net.ParseIP(argsWithProg[4])
	if dstIP == nil {
		fmt.Printf("Invalid Dest IP \n")
		return
	}

	dstPortVal, err := strconv.ParseUint(argsWithProg[5], 10, 16)
	if err != nil {
		fmt.Printf("Invalid Dest Port %v \n", err)
		return
	}

	srcPort := make([]byte, 2)
	binary.BigEndian.PutUint16(srcPort, uint16(srcPortVal))

	dstPort := make([]byte, 2)
	binary.BigEndian.PutUint16(dstPort, uint16(dstPortVal))

	if controlType == 1 {
		control_packet.SendAddMapping(srcIP[12:16], dstIP[12:16], srcPort, dstPort)
	}
	if controlType == 3 {
		control_packet.SendAddDestMapping(srcIP[12:16], dstIP[12:16], srcPort, dstPort)
	}
}

func main() {
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("Invalid Number of Arguments")
		fmt.Println("Looking for: (control type)")
		fmt.Println("Looking for: (control type) (sourceIP) (sourcePort) (destinationIP) (destinationPort)")
		return
	}

	controlType, err := strconv.Atoi(argsWithProg[1])
	if err != nil {
		fmt.Println("Error parsing control type")
		return
	}

	if controlType == 1 || controlType == 3 {
		addMapping(argsWithProg, controlType)
	} else if controlType == 2 {
		control_packet.SendListMappings()
	}
}
