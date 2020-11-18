package process_packet

import (
	"encoding/binary"
)

func pseudoHeader(srcIP [4]byte, dstIP [4]byte) uint32 {
	var csum uint32
	csum += (uint32(srcIP[0]) + uint32(srcIP[2])) << 8
	csum += uint32(srcIP[1]) + uint32(srcIP[3])
	csum += (uint32(dstIP[0]) + uint32(dstIP[2])) << 8
	csum += uint32(dstIP[1]) + uint32(dstIP[3])
	return csum
}

func computeChecksum(data []byte, csum uint32) uint16 {
	length := len(data) - 1
	for i := 0; i < length; i += 2 {
		// For our test packet, doing this manually is about 25% faster
		// (740 ns vs. 1000ns) than doing it by calling binary.BigEndian.Uint16.
		csum += uint32(data[i]) << 8
		csum += uint32(data[i+1])
	}
	if len(data)%2 == 1 {
		csum += uint32(data[length]) << 8
	}
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff)
	}
	return ^uint16(csum)
}

func udpCheckSum(data []byte) {
	srcIP, dstIP, _ := GetSrcDstIP(data[14:])
	csum := pseudoHeader(srcIP, dstIP)

	endIPHeader := ((uint8(data[14]) & 0x0F) << 2) + 14
	length := uint32(binary.BigEndian.Uint16(data[endIPHeader+4 : endIPHeader+6]))

	data[endIPHeader+6] = 0
	data[endIPHeader+7] = 0

	csum += uint32(data[23])
	csum += length & 0xffff
	csum += length >> 16

	checksum := computeChecksum(data[endIPHeader:len(data)], csum)

	checksumByte := make([]byte, 2)
	binary.BigEndian.PutUint16(checksumByte, checksum)

	data[endIPHeader+6] = checksumByte[0]
	data[endIPHeader+7] = checksumByte[1]
}

func tcpCheckSum(data []byte) {
	srcIP, dstIP, _ := GetSrcDstIP(data[14:])
	csum := pseudoHeader(srcIP, dstIP)

	endEthHeader := 14

	ihl := (uint8(data[14]) & 0x0F) << 2
	totalLength := binary.BigEndian.Uint16(data[16:18])
	length := uint32(totalLength) - uint32(ihl)

	endIPHeader := int(ihl) + endEthHeader

	data[endIPHeader+16] = 0
	data[endIPHeader+17] = 0

	csum += uint32(data[endEthHeader+9])
	csum += length & 0xffff
	csum += length >> 16

	checksum := computeChecksum(data[endIPHeader:len(data)], csum)

	//Printf("%#x \n", checksum)

	checksumByte := make([]byte, 2)
	binary.BigEndian.PutUint16(checksumByte, checksum)

	data[endIPHeader+16] = checksumByte[0]
	data[endIPHeader+17] = checksumByte[1]
}

func updateCheckSum(data []byte) {
	protocol := uint32(data[23])
	if protocol == 17 {
		udpCheckSum(data)
		ipCheckSum(data)
	} else {
		tcpCheckSum(data)
		ipCheckSum(data)
	}
}

func computeIPCheckSum(bytes []byte) uint16 {
	// Clear checksum bytes
	bytes[10] = 0
	bytes[11] = 0

	// Compute checksum
	var csum uint32
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8
		csum += uint32(bytes[i+1])
	}
	for {
		// Break when sum is less or equals to 0xFFFF
		if csum <= 65535 {
			break
		}
		// Add carry to the sum
		csum = (csum >> 16) + uint32(uint16(csum))
	}
	// Flip all the bits
	return ^uint16(csum)
}

func ipCheckSum(data []byte) {
	endEthHeader := 14
	endIPHeader := ((uint8(data[14]) & 0x0F) * 4) + 14
	checksum := computeIPCheckSum(data[endEthHeader:endIPHeader])

	checksumByte := make([]byte, 2)
	binary.BigEndian.PutUint16(checksumByte, checksum)

	data[endEthHeader+10] = checksumByte[0]
	data[endEthHeader+11] = checksumByte[1]
}
