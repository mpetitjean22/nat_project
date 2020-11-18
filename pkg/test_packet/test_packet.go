package test_packet

import (
	"log"
	"math/rand"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	err     error
	buffer  gopacket.SerializeBuffer
	options gopacket.SerializeOptions = gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
)

func CreateTestPacket(payload []byte) []byte {
	/*ethernetLayer := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0x08, 0x00, 0x27, 0xFD, 0x06, 0x32},
		DstMAC:       net.HardwareAddr{0x52, 0x54, 0x00, 0x12, 0x35, 0x02},
		EthernetType: layers.EthernetTypeIPv4,
	}*/
	ipLayer := &layers.IPv4{
		//SrcIP: net.IP{10, 0, 2, 15},
		//DstIP: net.IP{8, 8, 8, 8},

		SrcIP:    net.IP{10, 1, 0, 10},
		DstIP:    net.IP{1, 1, 1, 1},
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
		Id:       uint16(rand.Uint32()),
	}
	/*udpLayer := &layers.UDP{
		SrcPort: layers.UDPPort(55001),
		DstPort: layers.UDPPort(80),
	}
	udpLayer.SetNetworkLayerForChecksum(ipLayer) */
	tcpLayer := &layers.TCP{
		SrcPort: layers.TCPPort(80),
		DstPort: layers.TCPPort(80),
		Seq:     0,
		Window:  65535,
		SYN:     true,
	}
	//udpLayer.SetNetworkLayerForChecksum(ipLayer)
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	// And create the packet with the layers
	buffer = gopacket.NewSerializeBuffer()
	err := gopacket.SerializeLayers(buffer, options,
		//ethernetLayer,
		ipLayer,
		//udpLayer,
		tcpLayer,
		gopacket.Payload(payload),
	)
	if err != nil {
		log.Fatalf("marie should have checked the error return value omg it returns %v\n", err)
	}

	outgoingPacket := buffer.Bytes()
	return outgoingPacket
}
