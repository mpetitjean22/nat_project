/* This file contains some copy function in order to
 * be more FPGA friendly
 */

package process_packet

func TwoByteCopy(src []byte, startIdx int) [2]byte {
	res := [2]byte{}
	res[0] = src[startIdx]
	res[1] = src[startIdx+1]
	return res
}

func FourByteCopy(src []byte, startIdx int) [4]byte {
	res := [4]byte{}
	res[0] = src[startIdx]
	res[1] = src[startIdx+1]
	res[2] = src[startIdx+2]
	res[3] = src[startIdx+3]
	return res
}

func packetCopy(newPacket [65535]byte, newPacketStart int, data []byte, dataStart int, copyLength int) [65535]byte {
	var i int

	for i = 0; i < copyLength; i++ {
		newPacket[newPacketStart+i] = data[dataStart+i]
	}
	return newPacket
}
