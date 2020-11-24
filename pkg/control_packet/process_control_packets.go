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

// ProcessControlPacket processes the control packets where IP/Port is specified in Config.yaml
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
