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
	device       string = "tun2"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
	buffer       gopacket.SerializeBuffer
	options      gopacket.SerializeOptions = gopacket.SerializeOptions{FixLengths: true}
)

func createControlPacket(payload []byte) []byte {
	ipLayer := &layers.IPv4{
		SrcIP:    net.IP{10, 0, 0, 2}, // this should be adjusted
		DstIP:    net.IP{8, 8, 8, 8},
		Version:  4,
		TTL:      10,
		Protocol: layers.IPProtocolUDP,
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(80), // this should be adjusted
		DstPort: layers.UDPPort(80),
	}
	// And create the packet with the layers
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
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
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

// SendAddMapping Payload -> 1 byte: Control Type
//					0x01: add mapping
// 		   -> 4 bytes: Source IP
//		   -> 4 bytes: Destination IP
// 		   -> 2 bytes: Source Port
//         -> 2 bytes: Destination Port
func SendAddMapping(srcIP []byte, dstIP []byte, srcPort []byte, dstPort []byte) {
	payload := []byte{0x01}
	payload = append(payload, srcIP...)
	payload = append(payload, dstIP...)
	payload = append(payload, srcPort...)
	payload = append(payload, dstPort...)

	sendContolPacket(payload)
}

// SendAddDestMapping Payload -> 1 byte: Control Type
//					0x03: add mapping
// 		   -> 4 bytes: Source IP
//		   -> 4 bytes: Destination IP
// 		   -> 2 bytes: Source Port
//         -> 2 bytes: Destination Port
func SendAddDestMapping(srcIP []byte, dstIP []byte, srcPort []byte, dstPort []byte) {
	payload := []byte{0x03}
	payload = append(payload, srcIP...)
	payload = append(payload, dstIP...)
	payload = append(payload, srcPort...)
	payload = append(payload, dstPort...)

	sendContolPacket(payload)
}

// SendListMappings Payload -> 1 byte: Control Type
//					0x02: print mappings
func SendListMappings() {
	payload := []byte{0x02}
	sendContolPacket(payload)
}
