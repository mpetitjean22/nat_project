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
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
	buffer       gopacket.SerializeBuffer
	options      gopacket.SerializeOptions = gopacket.SerializeOptions{FixLengths: true}
)

func createControlPacket(payload []byte) []byte {
	ethernetLayer := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0xF0, 0x18, 0x98, 0x28, 0x0D, 0x06},
		DstMAC:       net.HardwareAddr{0xF0, 0x18, 0x98, 0x28, 0x0D, 0x06},
		EthernetType: layers.EthernetTypeIPv4,
	}
	ipLayer := &layers.IPv4{
		SrcIP:    net.IP{127, 0, 0, 1},
		DstIP:    net.IP{8, 8, 8, 8},
		Version:  4,
		TTL:      10,
		Protocol: layers.IPProtocolUDP,
	}
	udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(4320),
		DstPort: layers.UDPPort(80),
	}
	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()
	gopacket.SerializeLayers(buffer, options,
		ethernetLayer,
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

// Payload -> 1 byte: Control Type
//					0x01: add mapping
//					0x02: print mappings
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

func SendListMappings() {
	payload := []byte{0x02}
	sendContolPacket(payload)
}
