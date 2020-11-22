/* This file contains functions which modify the properties of
 * packets in order to create new packets. For example, rewriting
 * the source IP or rewriting the destination IP.
 */

package process_packet

import "fmt"

// WriteSource is kind of a hot mess right now but will cleanup
/* Functions below assume that data is a valid packet with IPv4 on top of UDP/TCP */
func WriteSource(data []byte, srcIP [4]byte, srcPort [2]byte) ([65535]byte, error) {
	var version byte
	version = data[0] >> 4

	if len(data) > 65535 {
		// for debugging but also just in case
		return [65535]byte{}, fmt.Errorf("Packet too large for buffer")
	}

	if version == 4 {
		var endEthHeader, endIPHeader, endIPEthHeaders int
		var newPacket [65535]byte

		endEthHeader = 14
		endIPHeader = int((uint8(data[0]) & 0x0F) * 4)
		endIPEthHeaders = endEthHeader + int(endIPHeader)

		// copy eth header
		ethHeader := []byte{0x52, 0x54, 0x00, 0x12, 0x35, 0x02, 0x08, 0x00, 0x27, 0xfd, 0x06, 0x32, 0x08, 0x00}
		copy(newPacket[0:14], ethHeader)

		// copy ipv4 header (with new source IP)
		copy(newPacket[14:26], data[0:12])
		copy(newPacket[26:30], srcIP[:])
		copy(newPacket[30:endIPEthHeaders], data[16:endIPHeader])

		// copy tcp/udp header (with new src port)
		//copy(newPacket[endIPEthHeaders:endIPEthHeaders+2], srcPort[:])

		// copy rest of packet
		copy(newPacket[endIPEthHeaders:], data[endIPHeader:])

		//fmt.Printf("%#v\n", newPacket[:14+len(data)])

		updateCheckSum(newPacket[:14+len(data)])
		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}

// WriteDestination is kind of a hot mess right now but will cleanup
func WriteDestination(data []byte, dstIP [4]byte, dstPort [2]byte) ([65535]byte, error) {
	var version byte

	version = data[14] >> 4

	if version == 4 {
		var endEthHeader, endIPHeader, endIPEthHeaders int
		var newPacket [65535]byte
		endEthHeader = 14
		endIPHeader = int((uint8(data[14]) & 0x0F) * 4)
		endIPEthHeaders = endEthHeader + int(endIPHeader)

		// copy eth header
		newPacket = packetCopy(newPacket, 0, data, 0, 14)

		// copy ipv4 header (with new dest IP)
		newPacket = packetCopy(newPacket, 14, data, 14, 16)
		newPacket = packetCopy(newPacket, 30, dstIP[:], 0, 4)

		newPacket = packetCopy(newPacket, 34, data, 34, endIPEthHeaders-34)

		// copy tcp/udp header (with new dest port)
		newPacket = packetCopy(newPacket, endIPEthHeaders, data, endIPEthHeaders, 4)
		//newPacket = packet_copy(newPacket, endIPEthHeaders+2, dstPort[:], 0, 2)

		// copy rest of packet
		newPacket = packetCopy(newPacket, endIPEthHeaders+4, data, endIPEthHeaders+4, len(data)-(endIPEthHeaders+4))
		updateCheckSum(newPacket[:len(data)])
		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}
