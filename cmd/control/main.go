package main

import (
	"encoding/binary"
	"fmt"
	"nat_project/pkg/control_packet"
	"net"
	"os"
	"strconv"
)

/* Create and Send various control packets to the inbound and outbound nat tables
 * via a command line interface.
 */
func main() {
	argsWithProg := os.Args
	if len(argsWithProg) < 2 {
		fmt.Println("Error: Invalid Number of Arguments")
		printControlOptions()
		return
	}

	controlType, err := strconv.Atoi(argsWithProg[1])
	if err != nil {
		fmt.Printf("Error: Could not parse control type given was: %v \n", argsWithProg[1])
		printControlOptions()
		return
	}

	if controlType == 1 || controlType == 3 {
		addMapping(argsWithProg, controlType)

	} else if controlType == 2 {
		control_packet.SendListMappings()
	} else {
		fmt.Printf("Error: Invalid control type given was %v \n", controlType)
		printControlOptions()
	}
}

func parseMappingArgs(argsWithProg []string) ([]byte, []byte, []byte, []byte, error) {
	if len(argsWithProg) != 6 {
		fmt.Println("Invalid Number of Arguments")
		return nil, nil, nil, nil, fmt.Errorf("Invalid Args")
	}

	fromIP := net.ParseIP(argsWithProg[2])
	if fromIP == nil {
		fmt.Printf("Error: fromIP given was invalid %v \n", argsWithProg[2])
		return nil, nil, nil, nil, fmt.Errorf("Invalid Args")
	}

	fromPortVal, err := strconv.ParseUint(argsWithProg[3], 10, 16)
	if err != nil {
		fmt.Printf("Error: fromPort given was invalid %v \n", argsWithProg[3])
		return nil, nil, nil, nil, fmt.Errorf("Invalid Args")
	}

	toIP := net.ParseIP(argsWithProg[4])
	if toIP == nil {
		fmt.Printf("Error: toIP given was invalid %v \n", argsWithProg[4])
		return nil, nil, nil, nil, fmt.Errorf("Invalid Args")
	}

	toPortVal, err := strconv.ParseUint(argsWithProg[5], 10, 16)
	if err != nil {
		fmt.Printf("Error: toPort given was invalid %v \n", argsWithProg[5])
		return nil, nil, nil, nil, fmt.Errorf("Invalid Args")
	}

	fromPort := make([]byte, 2)
	binary.BigEndian.PutUint16(fromPort, uint16(fromPortVal))

	toPort := make([]byte, 2)
	binary.BigEndian.PutUint16(toPort, uint16(toPortVal))

	return fromIP[12:16], toIP[12:16], fromPort, toPort, nil
}

func addMapping(argsWithProg []string, controlType int) {
	fromIP, fromPort, toIP, toPort, err := parseMappingArgs(argsWithProg)
	if err != nil {
		printControlOptions()
		return
	}
	if controlType == 1 {
		control_packet.SendAddMapping(fromIP, toIP, fromPort, toPort)
	}
	if controlType == 3 {
		control_packet.SendAddDestMapping(fromIP, toIP, fromPort, toPort)
	}
}

func printControlOptions() {
	fmt.Println("Looking for: (control type)")
	fmt.Println("   control types: ")
	fmt.Println("       2 -> list current mappings")
	fmt.Println("Looking for: (control type) (fromIP) (fromPort) (toIP) (toPort)")
	fmt.Println("   control types: ")
	fmt.Println("       1 -> add outbound mapping")
	fmt.Println("       3 -> add inbound mapping")
}
