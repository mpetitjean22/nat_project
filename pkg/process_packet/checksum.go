/* This file is dedicated towards computer checksums for
 * the udp/tcp header and also for the ip header. Allows us to
 * send well formed packets that will not be dropped.
 */

package process_packet;

import ("encoding/binary"); 

// pseudoHeader computes the checksum of the header
func pseudoHeader(srcIP [4]byte, dstIP [4]byte) uint32 {
	var csum uint32;
	csum += (uint32(srcIP[0]) + uint32(srcIP[2])) << 8;
	csum += uint32(srcIP[1]) + uint32(srcIP[3]);
	csum += (uint32(dstIP[0]) + uint32(dstIP[2])) << 8;
	csum += uint32(dstIP[1]) + uint32(dstIP[3]);
	return csum;
};

// computeChecksum uses the psuedo header checksum in order
// to determine the checksum of the tcp/udp packet
func computeChecksum(data []byte, csum uint32) uint16 {
	length := len(data) - 1;
	for i := 0; i < length; i += 2 {
		csum += uint32(data[i]) << 8;
		csum += uint32(data[i+1]);
	};
	if len(data)%2 == 1 {
		csum += uint32(data[length]) << 8;
	};
	for csum > 0xffff {
		csum = (csum >> 16) + (csum & 0xffff);
	};
	return ^uint16(csum);
};

// udpCheckSum computes the checksum of a udp packet
func udpCheckSum(data []byte) {
	srcIP, dstIP, _ := GetSrcDstIP(data[14:]);
	csum := pseudoHeader(srcIP, dstIP);

	endIPHeader := ((uint8(data[14]) & 0x0F) << 2) + 14;
	length := uint32(binary.BigEndian.Uint16(data[endIPHeader+4 : endIPHeader+6]));

	data[endIPHeader+6] = 0;
	data[endIPHeader+7] = 0;

	csum += uint32(data[23]);
	csum += length & 0xffff;
	csum += length >> 16;

	checksum := computeChecksum(data[endIPHeader:len(data)], csum);

	checksumByte := make([]byte, 2);
	binary.BigEndian.PutUint16(checksumByte, checksum);

	data[endIPHeader+6] = checksumByte[0];
	data[endIPHeader+7] = checksumByte[1];
};

// tcpCheckSum computes the checksum of tcp packet
func tcpCheckSum(data []byte) {
	srcIP, dstIP, _ := GetSrcDstIP(data[14:]);
	csum := pseudoHeader(srcIP, dstIP);

	endEthHeader := 14;

	ihl := (uint8(data[14]) & 0x0F) << 2;
	totalLength := binary.BigEndian.Uint16(data[16:18]);
	length := uint32(totalLength) - uint32(ihl);

	endIPHeader := int(ihl) + endEthHeader;

	data[endIPHeader+16] = 0;
	data[endIPHeader+17] = 0;

	csum += uint32(data[endEthHeader+9]);
	csum += length & 0xffff;
	csum += length >> 16;

	checksum := computeChecksum(data[endIPHeader:len(data)], csum);
	checksumByte := make([]byte, 2);
	binary.BigEndian.PutUint16(checksumByte, checksum);

	data[endIPHeader+16] = checksumByte[0];
	data[endIPHeader+17] = checksumByte[1];
};

// updateCheckSum determines if packet is a udp or tcp packet
// and computes the appropriate checksum
func updateCheckSum(data []byte) {
	protocol := uint32(data[23]);
	if protocol == 17 {
		udpCheckSum(data);
		ipCheckSum(data);
	} else {
		tcpCheckSum(data);
		ipCheckSum(data);
	};
};

// computeIPCheckSum computes the checksum of the IP header
func computeIPCheckSum(bytes []byte) uint16 {
	bytes[10] = 0;
	bytes[11] = 0;

	var csum uint32;
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8;
		csum += uint32(bytes[i+1]);
	};
	for {
		if csum <= 65535 {
			break;
		};
		csum = (csum >> 16) + uint32(uint16(csum));
	};
	return ^uint16(csum);
};

// ipCheckSum finds the checksum of the ip header
func ipCheckSum(data []byte) {
	endEthHeader := 14;
	endIPHeader := ((uint8(data[14]) & 0x0F) * 4) + 14;
	checksum := computeIPCheckSum(data[endEthHeader:endIPHeader]);

	checksumByte := make([]byte, 2);
	binary.BigEndian.PutUint16(checksumByte, checksum);

	data[endEthHeader+10] = checksumByte[0];
	data[endEthHeader+11] = checksumByte[1];
};
