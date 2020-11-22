/* This file contains functions which parse packets received
 * as control packets and performed the necessary operations
 * that they indicate.
 */

package control_packet

import (
	"fmt"
	"nat_project/pkg/nat"
	"nat_project/pkg/process_packet"
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
		srcIP := process_packet.FourByteCopy(payload, 1)
		srcPort := process_packet.TwoByteCopy(payload, 5)

		dstIP := process_packet.FourByteCopy(payload, 7)
		dstPort := process_packet.TwoByteCopy(payload, 11)

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
