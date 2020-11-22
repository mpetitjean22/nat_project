package control_packet

import (
	"log"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	device      string        = "tun2" // TODO: make generalizable!
	snapshotLen int32         = 1024
	promiscuous bool          = false
	timeout     time.Duration = 30 * time.Second

	buffer  gopacket.SerializeBuffer
	options gopacket.SerializeOptions = gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
)

// SendAddMapping creates and sends a packet with the following payload:
//			-> 1 byte: Control Type
//				0x01: add mapping
// 		   	-> 4 bytes: from IP
//		   	-> 4 bytes: to IP
// 		   	-> 2 bytes: from Port
//         	-> 2 bytes: to Port
func SendAddMapping(srcIP []byte, dstIP []byte, srcPort []byte, dstPort []byte) {
	payload := []byte{0x01}
	payload = append(payload, srcIP...)
	payload = append(payload, dstIP...)
	payload = append(payload, srcPort...)
	payload = append(payload, dstPort...)

	sendContolPacket(payload)
}

// SendAddDestMapping creates and sends a packet with the following payload:
// 			-> 1 byte: Control Type
//				0x03: add mapping
// 			-> 4 bytes: from IP
//		   	-> 4 bytes: to IP
// 		   	-> 2 bytes: from Port
//         	-> 2 bytes: to Port
func SendAddDestMapping(srcIP []byte, dstIP []byte, srcPort []byte, dstPort []byte) {
	payload := []byte{0x03}
	payload = append(payload, srcIP...)
	payload = append(payload, dstIP...)
	payload = append(payload, srcPort...)
	payload = append(payload, dstPort...)

	sendContolPacket(payload)
}

// SendListMappings creates and sends a packet with the following payload:
//			-> 1 byte: Control Type
//				0x02: print mappings
func SendListMappings() {
	payload := []byte{0x02}
	sendContolPacket(payload)
}

func createControlPacket(payload []byte) []byte {
	ipLayer := &layers.IPv4{
		SrcIP:    net.IP{10, 0, 0, 2}, // TODO: make generalizable!
		DstIP:    net.IP{8, 8, 8, 8},  // TODO: make generalizable!
		Version:  4,
		TTL:      10,
		Protocol: layers.IPProtocolUDP,
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(80), // TODO: make generalizable!
		DstPort: layers.UDPPort(80), // TODO: make generalizable!
	}
	udpLayer.SetNetworkLayerForChecksum(ipLayer)

	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		ipLayer,
		udpLayer,
		gopacket.Payload(payload),
	)
	outgoingPacket := buffer.Bytes()
	return outgoingPacket
}

func sendContolPacket(payload []byte) {
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	outgoingPacket := createControlPacket(payload)
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal(err)
	}
}
