package control_packet

import (
	"nat_project/pkg/nat"
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/stretchr/testify/assert"
)

func createGoPacket(data []byte) gopacket.Packet {
	return gopacket.NewPacket(data, layers.LayerTypeIPv4, gopacket.Default)
}

func getGoPacketValues(packet gopacket.Packet, t *testing.T) (dstIP net.IP, dstPort uint16) {
	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ip4Layer != nil {
		ip, _ := ip4Layer.(*layers.IPv4)
		dstIP = ip.DstIP
	} else {
		t.Errorf("Not IPv4")
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		dstPort = uint16(udp.DstPort)
	} else {
		t.Errorf("Not UDP")
	}

	return
}

func TestAddMappingPacket(t *testing.T) {
	payload := []byte{0x01}
	rawPacket := createControlPacket(payload)

	goPacket := createGoPacket(rawPacket)
	expIP, expPort := getGoPacketValues(goPacket, t)

	assert.Equal(t, expIP, net.IP(nat.Configs.Ctrl.IP[:]))
	assert.Equal(t, expPort, nat.Configs.Ctrl.Port)
}
