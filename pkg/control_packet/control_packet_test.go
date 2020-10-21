package control_packet

import (
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

func createGoPacket(data []byte) gopacket.Packet {
	return gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
}

func getGoPacketValues(packet gopacket.Packet, t *testing.T) (dstIP net.IP, dstPort uint16) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer == nil {
		t.Errorf("No Ethernet Type")
	}

	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	ip6Layer := packet.Layer(layers.LayerTypeIPv6)

	if ip4Layer != nil {
		ip, _ := ip4Layer.(*layers.IPv4)
		dstIP = ip.DstIP
	} else if ip6Layer != nil {
		ip, _ := ip6Layer.(*layers.IPv6)
		dstIP = ip.DstIP
	} else {
		t.Errorf("Not IPv4 or IPv6")
	}

	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		dstPort = uint16(tcp.DstPort)
	} else if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		dstPort = uint16(udp.DstPort)
	} else {
		t.Errorf("Not UDP or TCP")
	}

	return
}

func TestAddMappingPacket(t *testing.T) {
	payload := []byte{0x01}
	rawPacket := createControlPacket(payload)

	goPacket := createGoPacket(rawPacket)
	expIP, expPort := getGoPacketValues(goPacket, t)

	assert.Equal(t, expIP, net.IP{0x08, 0x08, 0x08, 0x08})
	assert.Equal(t, expPort, uint16(80))
}
