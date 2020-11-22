package control_packet

import (
	"fmt"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
)

// Destination IP Address/Port for Control Packets
var (
	ControlIP   = [4]byte{0x08, 0x08, 0x08, 0x08} // TODO: make generalizable!
	ControlPort = [2]byte{0x00, 0x50}             // TODO: make generalizable!
)

// ProcessControlPacket processes the control packets with the following form:
// 			-> dstIP: 8.8.8.8
// 			-> dstPort: 80
func ProcessControlPacket(packet []byte, outboundNat *nat.Table, inboundNat *nat.Table) {
	ihl := uint8(packet[0]) & 0x0F
	payload := packet[8+(ihl*4):]
	controlType := payload[0]

	if controlType == 1 || controlType == 3 {
		srcIP := get_packets.Four_byte_copy(payload, 1)
		srcPort := get_packets.Two_byte_copy(payload, 5)

		dstIP := get_packets.Four_byte_copy(payload, 7)
		dstPort := get_packets.Two_byte_copy(payload, 11)

		if controlType == 1 {
			outboundNat.AddMapping(srcIP, srcPort, dstIP, dstPort)
		} else {
			inboundNat.AddMapping(srcIP, srcPort, dstIP, dstPort)
		}
	} else if controlType == 2 {
		fmt.Println("Outbound")
		outboundNat.PrettyPrintTable()

		fmt.Println("Inbound")
		inboundNat.PrettyPrintTable()
	}
}
