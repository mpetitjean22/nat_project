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
