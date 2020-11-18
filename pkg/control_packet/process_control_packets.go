package control_packet

import (
	"fmt"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
)

// Destination IP Address/Port for Control Packets
var (
	ControlIP   = [4]byte{0x08, 0x08, 0x08, 0x08}
	ControlPort = [2]byte{0x00, 0x50}
)

func pp_table(mappings map[nat.IPAddress]*nat.IPAddress) {
	fmt.Println("--------------------------")
	for key, value := range mappings {
		fmt.Printf("%v to %v \n", key, *value)
	}
	fmt.Println("--------------------------")
}

// Decodes a control packet
// We will define control packets as having (but can be changed above):
// 			-> dstIP: 8.8.8.8
// 			-> dstPort: 80
func ProcessControlPacket(packet []byte, outbound_nat *nat.NAT_Table, inbound_nat *nat.NAT_Table) {
	ihl := uint8(packet[0]) & 0x0F
	payload := packet[8+(ihl*4):]
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
		pp_table(outbound_nat.ListMappings())

		fmt.Println("Inbound")
		pp_table(inbound_nat.ListMappings())
	}
}
