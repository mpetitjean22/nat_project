package process_packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"nat_project/pkg/get_packets"
	"reflect"
)

func getSrcDstPortIPv4(data []byte) ([2]byte, [2]byte, error) {
	var length uint16
	var ihl uint8
	var protocol byte
	var payload []byte
	var sPort, dPort [2]byte

	length = binary.BigEndian.Uint16(data[2:4])
	protocol = data[9]
	if protocol != 6 && protocol != 17 {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Not TCP or UDP")
	}

	ihl = uint8(data[0]) & 0x0F
	if length < 20 {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Invalid (too small) IP length (%d < 20)", length)
	} else if ihl < 5 {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Invalid (too small) IP header length (%d < 5)", ihl)
	} else if int(ihl*4) > int(length) {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Invalid IP header length > IP length (%d > %d)", ihl, length)
	}

	payload = data[ihl*4:]

	sPort = get_packets.Two_byte_copy(payload, 0)
	dPort = get_packets.Two_byte_copy(payload, 2)

	return sPort, dPort, nil
}

// Dead code for now -- skelaton for IPv6 support
/* func getSrcDstPortIPv6(data []byte) (uint16, uint16, error) {
	protocol := data[6]
	if protocol != 6 && protocol != 17 {
		return 0, 0, fmt.Errorf("Not TCP or UDP")
	}

	// TODO: Implement extracting Source and Dest Ports
	// from the the payload with IPv6 header (having some trouble
	// figuring out how big the IPv6 head is)
	return 0, 0, nil
} */

func GetSrcDstPort(data []byte) ([2]byte, [2]byte, error) {
	var version byte
	version = data[0] >> 4

	if version == 4 {
		return getSrcDstPortIPv4(data)
	}
	/* Revisit IPv6 Support!
	else if version == 6 {
		return getSrcDstPortIPv6(data)
	}
	*/

	return [2]byte{}, [2]byte{}, nil
}

func GetSrcDstIP(data []byte) ([4]byte, [4]byte, error) {
	var version byte
	var srcIP, dstIP [4]byte

	if len(data) < 20 {
		return [4]byte{}, [4]byte{}, fmt.Errorf("Invalid ip4 header. Length %d less than 20", len(data))
	}
	version = data[0] >> 4
	if version == 4 {
		srcIP = get_packets.Four_byte_copy(data, 12)
		dstIP = get_packets.Four_byte_copy(data, 16)
		return srcIP, dstIP, nil
	}

	/* Will have to revist the IPv6 Support (since it needs more than 4 bytes)
		else if version == 6 {
		srcIP := net.IP(data[8:24])
		dstIP := net.IP(data[24:40])
		return srcIP, dstIP, nil
	}
	*/

	return [4]byte{}, [4]byte{}, fmt.Errorf("Not Valid Version")
}

func GetMacAddress(data []byte) bool {
	var Marie_MAC, dstMac, srcMac []byte

	Marie_MAC = []byte{0xF0, 0x18, 0x98, 0x28, 0x0D, 0x06}
	dstMac = data[:6]
	srcMac = data[6:12]

	if reflect.DeepEqual(Marie_MAC, dstMac) {
		return true
	}
	if reflect.DeepEqual(Marie_MAC, srcMac) {
		return false
	}
	return false
}

func GetEthProtocol(data []byte) (uint16, error) {
	var ethernetType uint16

	if len(data) < 14 {
		return 0, errors.New("Ethernet packet too small")
	}

	ethernetType = binary.BigEndian.Uint16(data[12:14])
	return ethernetType, nil
}

func packet_copy(newPacket [65535]byte, newPacketStart int, data []byte, dataStart int, copyLength int) [65535]byte {
	var i int

	for i = 0; i < copyLength; i++ {
		newPacket[newPacketStart+i] = data[dataStart+i]
	}
	return newPacket
}

/* Functions below assume that data is a valid packet with IPv4 on top of UDP/TCP */
func WriteSource(data []byte, srcIP [4]byte, srcPort [2]byte) ([65535]byte, error) {
	var version byte
	version = data[14] >> 4

	if len(data) > 65535 {
		// for debugging but also just in case
		return [65535]byte{}, fmt.Errorf("Packet too large for buffer")
	}

	if version == 4 {
		var endEthHeader, endIPHeader, endIPEthHeaders int
		var newPacket [65535]byte

		endEthHeader = 14
		endIPHeader = int((uint8(data[14]) & 0x0F) * 4)
		endIPEthHeaders = endEthHeader + int(endIPHeader)

		// copy eth header
		newPacket = packet_copy(newPacket, 0, data, 0, 14)

		// copy ipv4 header (with new source IP)
		newPacket = packet_copy(newPacket, 14, data, 14, 12)
		newPacket = packet_copy(newPacket, 26, srcIP[:], 0, 4)
		newPacket = packet_copy(newPacket, 30, data, 30, endIPEthHeaders-30)

		// copy tcp/udp header (with new src port)
		newPacket = packet_copy(newPacket, endIPEthHeaders, srcPort[:], 0, 2)

		// copy rest of packet
		newPacket = packet_copy(newPacket, endIPEthHeaders+2, data, endIPEthHeaders+2, len(data)-(endIPEthHeaders+2))

		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}

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
		newPacket = packet_copy(newPacket, 0, data, 0, 14)

		// copy ipv4 header (with new dest IP)
		newPacket = packet_copy(newPacket, 14, data, 14, 16)
		newPacket = packet_copy(newPacket, 30, dstIP[:], 0, 4)

		newPacket = packet_copy(newPacket, 34, data, 34, endIPEthHeaders-34)

		// copy tcp/udp header (with new dest port)
		newPacket = packet_copy(newPacket, endIPEthHeaders, data, endIPEthHeaders, 2)
		newPacket = packet_copy(newPacket, endIPEthHeaders+2, dstPort[:], 0, 2)

		// copy rest of packet
		newPacket = packet_copy(newPacket, endIPEthHeaders+4, data, endIPEthHeaders+4, len(data)-(endIPEthHeaders+4))
		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}
