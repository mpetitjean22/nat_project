/* This file contains function which parse raw packets in the
 * form of byte arrays in order to extract information about
 * them and their contents.
 */

package process_packet

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// returns the source and dest port of either udp or tcp header
// in the form of a fixed length array.
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

	sPort = TwoByteCopy(payload, 0)
	dPort = TwoByteCopy(payload, 2)

	return sPort, dPort, nil
}

// GetSrcDstPort determines if packet has an IPv4 header, and if
// so returns the source and dest ports in fixed length arrays. This
// code is in place in case IPv6 support is implemented.
func GetSrcDstPort(data []byte) ([2]byte, [2]byte, error) {
	var version byte
	version = data[0] >> 4

	if version == 4 {
		return getSrcDstPortIPv4(data)
	}

	return [2]byte{}, [2]byte{}, nil
}

// GetSrcDstIP gets the Source and Desintation IPs from an IPv4 header
func GetSrcDstIP(data []byte) ([4]byte, [4]byte, error) {
	var version byte
	var srcIP, dstIP [4]byte

	if len(data) < 20 {
		return [4]byte{}, [4]byte{}, fmt.Errorf("Invalid ip4 header. Length %d less than 20", len(data))
	}
	version = data[0] >> 4
	if version == 4 {
		srcIP = FourByteCopy(data, 12)
		dstIP = FourByteCopy(data, 16)
		return srcIP, dstIP, nil
	}

	return [4]byte{}, [4]byte{}, fmt.Errorf("Not Valid Version")
}

// GetEthProtocol returns the ethernet protocol of the packet
func GetEthProtocol(data []byte) (uint16, error) {
	var ethernetType uint16

	if len(data) < 14 {
		return 0, errors.New("Ethernet packet too small")
	}

	ethernetType = binary.BigEndian.Uint16(data[12:14])
	return ethernetType, nil
}
