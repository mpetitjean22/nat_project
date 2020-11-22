/* This file contains functions which modify the properties of
 * packets in order to create new packets. For example, rewriting
 * the source IP or rewriting the destination IP.
 */

package process_packet

import "fmt"

// WriteSource takes a packet (without a ethernet header!) and rewrites the source IP
// and source port. Recomputes the checksum and returns the new raw byte packet with ethernet
// header designed to be sent out on eth interface.
func WriteSource(data []byte, srcIP [4]byte, srcPort [2]byte) ([65535]byte, error) {
	var version byte
	version = data[0] >> 4

	if len(data) > 65535 {
		return [65535]byte{}, fmt.Errorf("Packet too large for buffer")
	}

	if version == 4 {
		var endEthHeader, endIPHeader, endIPEthHeaders int
		var newPacket [65535]byte

		endEthHeader = 14
		endIPHeader = int((uint8(data[0]) & 0x0F) * 4)
		endIPEthHeaders = endEthHeader + int(endIPHeader)

		// copy eth header
		// TODO: make generalizable!
		ethHeader := []byte{0x52, 0x54, 0x00, 0x12, 0x35, 0x02, 0x08, 0x00, 0x27, 0xfd, 0x06, 0x32, 0x08, 0x00}
		newPacket = packetCopy(newPacket, 0, ethHeader, 0, 14)

		// copy ipv4 header (with new source IP)
		newPacket = packetCopy(newPacket, 14, data, 0, 12)
		newPacket = packetCopy(newPacket, 26, srcIP[:], 0, 4)
		newPacket = packetCopy(newPacket, 30, data, 16, endIPHeader-16)

		// copy src port to tcp/udp header
		newPacket = packetCopy(newPacket, endIPEthHeaders, srcPort[:], 0, 2)

		// copy rest of packet
		newPacket = packetCopy(newPacket, endIPEthHeaders+2, data, endIPHeader+2, len(data)-(endIPHeader+2))

		updateCheckSum(newPacket[:14+len(data)])
		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}

// WriteDestination takes a packets (with ethernet header!) and rewritest the destination IP
// and destination port. Recomputes the checksum and returns the new raw byte packets
// with ethernet header (which should be dropped when sending on tun interface)
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
		newPacket = packetCopy(newPacket, endIPEthHeaders, data, endIPEthHeaders, 2)
		newPacket = packetCopy(newPacket, endIPEthHeaders+2, dstPort[:], 0, 2)

		// copy rest of packet
		newPacket = packetCopy(newPacket, endIPEthHeaders+4, data, endIPEthHeaders+4, len(data)-(endIPEthHeaders+4))

		updateCheckSum(newPacket[:len(data)])
		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}
