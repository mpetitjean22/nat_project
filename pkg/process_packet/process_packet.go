package process_packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

func getSrcDstPortIPv4(data []byte) ([2]byte, [2]byte, error) {
	length := binary.BigEndian.Uint16(data[2:4])
	protocol := data[9]
	if protocol != 6 && protocol != 17 {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Not TCP or UDP")
	}

	ihl := uint8(data[0]) & 0x0F
	if length < 20 {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Invalid (too small) IP length (%d < 20)", length)
	} else if ihl < 5 {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Invalid (too small) IP header length (%d < 5)", ihl)
	} else if int(ihl*4) > int(length) {
		return [2]byte{}, [2]byte{}, fmt.Errorf("Invalid IP header length > IP length (%d > %d)", ihl, length)
	}

	payload := data[ihl*4:]

	sPort := [2]byte{}
	copy(sPort[:], payload[0:2])

	dPort := [2]byte{}
	copy(dPort[:], payload[2:4])

	return sPort, dPort, nil
}

func getSrcDstPortIPv6(data []byte) (uint16, uint16, error) {
	protocol := data[6]
	if protocol != 6 && protocol != 17 {
		return 0, 0, fmt.Errorf("Not TCP or UDP")
	}

	// TODO: Implement extracting Source and Dest Ports
	// from the the payload with IPv6 header (having some trouble
	// figuring out how big the IPv6 head is)
	return 0, 0, nil
}

func GetSrcDstPort(data []byte) ([2]byte, [2]byte, error) {
	version := data[0] >> 4

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
	if len(data) < 20 {
		return [4]byte{}, [4]byte{}, fmt.Errorf("Invalid ip4 header. Length %d less than 20", len(data))
	}
	version := data[0] >> 4
	if version == 4 {
		srcIP := [4]byte{}
		copy(srcIP[:], data[12:16])

		dstIP := [4]byte{}
		copy(dstIP[:], data[16:20])
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
	Marie_MAC := []byte{0xF0, 0x18, 0x98, 0x28, 0x0D, 0x06}
	dstMac := data[:6]
	srcMac := data[6:12]

	if reflect.DeepEqual(Marie_MAC, dstMac) {
		return true
	}
	if reflect.DeepEqual(Marie_MAC, srcMac) {
		return false
	}
	return false
}

func GetEthProtocol(data []byte) (uint16, error) {
	if len(data) < 14 {
		return 0, errors.New("Ethernet packet too small")
	}

	ethernetType := binary.BigEndian.Uint16(data[12:14])
	return ethernetType, nil
}

/* Functions below assume that data is a valid packet with IPv4 on top of UDP/TCP */
func WriteSource(data []byte, srcIP [4]byte, srcPort [2]byte) []byte {
	version := data[14] >> 4

	if version == 4 {
		endEthHeader := 14
		endIPHeader := (uint8(data[14]) & 0x0F) * 4
		endIPEthHeaders := endEthHeader + int(endIPHeader)

		newPacket := make([]byte, len(data))

		// copy eth header
		copy(newPacket[:14], data[:14])

		// copy ipv4 header (with new source IP)
		copy(newPacket[14:26], data[14:26])
		copy(newPacket[26:30], srcIP[:])
		copy(newPacket[30:endIPEthHeaders], data[30:endIPEthHeaders])

		// copy tcp/udp header (with new src port)
		copy(newPacket[endIPEthHeaders:endIPEthHeaders+2], srcPort[:])

		// copy rest of packet
		copy(newPacket[endIPEthHeaders+2:], data[endIPEthHeaders+2:])

		return newPacket
	}

	return nil
}

func WriteDestination(data []byte, dstIP [4]byte, dstPort [2]byte) []byte {
	version := data[14] >> 4

	if version == 4 {
		endEthHeader := 14
		endIPHeader := (uint8(data[14]) & 0x0F) * 4
		endIPEthHeaders := endEthHeader + int(endIPHeader)

		newPacket := make([]byte, len(data))

		// copy eth header
		copy(newPacket[:14], data[:14])

		// copy ipv4 header (with new dest IP)
		copy(newPacket[14:30], data[14:30])
		copy(newPacket[30:34], dstIP[:])

		copy(newPacket[34:endIPEthHeaders], data[34:endIPEthHeaders])

		// copy tcp/udp header (with new dest port)
		copy(newPacket[endIPEthHeaders:endIPEthHeaders+2], data[endIPEthHeaders:endIPEthHeaders+2])
		copy(newPacket[endIPEthHeaders+2:endIPEthHeaders+4], dstPort[:])

		// copy rest of packet
		copy(newPacket[endIPEthHeaders+4:], data[endIPEthHeaders+4:])
		return newPacket
	}

	return nil
}
