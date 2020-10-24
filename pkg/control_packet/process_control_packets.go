package control_packet

import (
	"fmt"
	"nat_project/pkg/nat"
)

// Decodes a control packet
// We will define control packets as having:
// 			-> dstIP: 8.8.8.8
// 			-> dstPort: 80
func ProcessControlPacket(packet []byte) {
	ihl := uint8(packet[14]) & 0x0F
	payload := packet[14+8+(ihl*4):]

	controlType := payload[0]

	srcIP := [4]byte{}
	copy(srcIP[:], payload[1:5])

	dstIP := [4]byte{}
	copy(dstIP[:], payload[5:9])

	srcPort := [2]byte{}
	copy(srcPort[:], payload[9:11])

	dstPort := [2]byte{}
	copy(dstPort[:], payload[11:13])

	if controlType == 1 {
		nat.AddMapping(srcIP, srcPort, dstIP, dstPort)
	} else if controlType == 2 {
		fmt.Println(nat.ListMappings())
	}
}
