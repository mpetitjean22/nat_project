package control_packet

import (
	"encoding/binary"
	"fmt"
	"nat_project/pkg/nat"
	"net"
	"strconv"
)

// Decodes a control packet
// We will define control packets as having:
// 			-> dstIP: 8.8.8.8
// 			-> dstPort: 80
func ProcessControlPacket(packet []byte) {
	ihl := uint8(packet[14]) & 0x0F
	payload := packet[14+8+(ihl*4):]

	controlType := payload[0]
	srcIP := net.IP(payload[1:5])
	dstIP := net.IP(payload[5:9])
	srcPort := binary.BigEndian.Uint16(payload[9:11])
	dstPort := binary.BigEndian.Uint16(payload[11:13])

	if controlType == 1 {
		nat.AddMapping(srcIP.String(), strconv.Itoa(int(srcPort)), dstIP.String(), strconv.Itoa(int(dstPort)))
	} else if controlType == 2 {
		fmt.Println(nat.ListMappings())
	}
}
