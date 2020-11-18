package process_packet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"nat_project/pkg/get_packets"
	"net"
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

func GetSrcDstPort(data []byte) ([2]byte, [2]byte, error) {
	var version byte
	version = data[0] >> 4

	if version == 4 {
		return getSrcDstPortIPv4(data)
	}

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

	return [4]byte{}, [4]byte{}, fmt.Errorf("Not Valid Version")
}

func IsInboundPacket(data []byte) bool {
	var macAddr, srcMac, dstMac []byte
	var netInterface *net.Interface
	var err error

	netInterface, err = net.InterfaceByName("enp0s3")
	if err != nil {
		return false
	}

	macAddr = netInterface.HardwareAddr
	dstMac = data[:6]
	srcMac = data[6:12]

	if reflect.DeepEqual(macAddr, dstMac) {
		return false
	}
	if reflect.DeepEqual(macAddr, srcMac) {
		return true
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

		// BEGIN REMOVE
		if bytes.Equal(newPacket[30:34], []byte{1, 2, 3, 4}) {
			copy(newPacket[30:34], []byte{10, 0, 2, 15})
		}
		// END REMOVE

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
		newPacket = packet_copy(newPacket, 0, data, 0, 14)

		// copy ipv4 header (with new dest IP)
		newPacket = packet_copy(newPacket, 14, data, 14, 16)
		newPacket = packet_copy(newPacket, 30, dstIP[:], 0, 4)

		// BEGIN REMOVE
		if bytes.Equal(data[26:30], []byte{10, 0, 2, 15}) {
			newPacket = packet_copy(newPacket, 26, []byte{1, 2, 3, 4}, 0, 4)
		}
		// END REMOVE

		newPacket = packet_copy(newPacket, 34, data, 34, endIPEthHeaders-34)

		// copy tcp/udp header (with new dest port)
		newPacket = packet_copy(newPacket, endIPEthHeaders, data, endIPEthHeaders, 4)
		//newPacket = packet_copy(newPacket, endIPEthHeaders+2, dstPort[:], 0, 2)

		// copy rest of packet
		newPacket = packet_copy(newPacket, endIPEthHeaders+4, data, endIPEthHeaders+4, len(data)-(endIPEthHeaders+4))
		updateCheckSum(newPacket[:len(data)])
		return newPacket, nil
	}

	return [65535]byte{}, fmt.Errorf("Invalid IP Version")
}
