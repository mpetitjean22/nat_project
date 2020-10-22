package process_packet

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func getSrcDstPortIPv4(data []byte) (uint16, uint16, error) {
	length := binary.BigEndian.Uint16(data[2:4])
	protocol := data[9]
	if protocol != 6 && protocol != 17 {
		return 0, 0, fmt.Errorf("Not TCP or UDP")
	}

	ihl := uint8(data[0]) & 0x0F
	if length < 20 {
		return 0, 0, fmt.Errorf("Invalid (too small) IP length (%d < 20)", length)
	} else if ihl < 5 {
		return 0, 0, fmt.Errorf("Invalid (too small) IP header length (%d < 5)", ihl)
	} else if int(ihl*4) > int(length) {
		return 0, 0, fmt.Errorf("Invalid IP header length > IP length (%d > %d)", ihl, length)
	}

	payload := data[ihl*4:]
	sPort := binary.BigEndian.Uint16(payload[0:2])
	dPort := binary.BigEndian.Uint16(payload[2:4])
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

func GetSrcDstPort(data []byte) (uint16, uint16, error) {
	version := data[0] >> 4

	if version == 4 {
		return getSrcDstPortIPv4(data)
	} else if version == 6 {
		return getSrcDstPortIPv6(data)
	}

	return 0, 0, nil
}

func GetSrcDstIP(data []byte) (net.IP, net.IP, error) {
	if len(data) < 20 {
		return nil, nil, fmt.Errorf("Invalid ip4 header. Length %d less than 20", len(data))
	}
	version := data[0] >> 4
	if version == 4 {
		srcIP := net.IP(data[12:16])
		dstIP := net.IP(data[16:20])
		return srcIP, dstIP, nil
	} else if version == 6 {
		srcIP := net.IP(data[8:24])
		dstIP := net.IP(data[24:40])
		return srcIP, dstIP, nil
	}

	return nil, nil, fmt.Errorf("Not Valid Version")
}

func GetEthProtocol(data []byte) (uint16, error) {
	if len(data) < 14 {
		return 0, errors.New("Ethernet packet too small")
	}

	ethernetType := binary.BigEndian.Uint16(data[12:14])
	return ethernetType, nil
}

/* Assumes that data is a valid packet with IPv4 on top of UDP/TCP */
func WriteDstIP(data []byte) []byte {
	version := data[14] >> 4

	if version == 4 {
		endEthHeader := 14
		endIPHeader := (uint8(data[14]) & 0x0F) * 4
		endIPEthHeaders := endEthHeader + int(endIPHeader)

		newPacket := make([]byte, len(data))

		sourceIP := []byte{0x02, 0x02, 0x02, 0x02} // hard coded for now
		sourcePort := []byte{0x00, 0x50}

		// copy eth header
		copy(newPacket[:14], data[:14])

		// copy ipv4 header (with new source IP)
		copy(newPacket[14:26], data[14:26])
		copy(newPacket[26:30], sourceIP)
		copy(newPacket[30:endIPEthHeaders], data[30:endIPEthHeaders])

		// copy tcp/udp header (with new dest port)
		copy(newPacket[endIPEthHeaders:endIPEthHeaders+2], sourcePort)

		// copy rest of packet
		copy(newPacket[endIPEthHeaders+2:], data[endIPEthHeaders+2:])

		return newPacket
	}

	return nil
}
