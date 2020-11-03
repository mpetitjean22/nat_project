package control_packet

import (
	"fmt"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
)

// Decodes a control packet
// We will define control packets as having:
// 			-> dstIP: 8.8.8.8
// 			-> dstPort: 80
func ProcessControlPacket(packet []byte, outbound_nat *nat.NAT_Table, inbound_nat *nat.NAT_Table) {
	ihl := uint8(packet[14]) & 0x0F
	payload := packet[14+8+(ihl*4):]
	controlType := payload[0]

	if controlType == 1 || controlType == 3 {
		srcIP := get_packets.Four_byte_copy(payload, 1)
		dstIP := get_packets.Four_byte_copy(payload, 5)

		srcPort := get_packets.Two_byte_copy(payload, 9)
		dstPort := get_packets.Two_byte_copy(payload, 11)

		if controlType == 1 {
			outbound_nat.AddMapping(srcIP, srcPort, dstIP, dstPort)
		} else {
			inbound_nat.AddMapping(srcIP, srcPort, dstIP, dstPort)
		}
	} else if controlType == 2 {
		fmt.Println("Outbound")
		fmt.Println(outbound_nat.ListMappings())

		fmt.Println("Inbound")
		fmt.Println(inbound_nat.ListMappings())
	}
}
