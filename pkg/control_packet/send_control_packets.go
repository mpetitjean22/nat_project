/* This file contains functions which contructs and sends
 * control packets in order to be used with the CLI tool.
 */

package control_packet;

import ("log"); 
import ("nat_project/pkg/nat"); 
import ("net"); 
import ("time"); 

import ("github.com/google/gopacket"); 
import ("github.com/google/gopacket/layers"); 
import ("github.com/google/gopacket/pcap"); 

// Parse does not support global variables 
// START REMOVE
/*var (
	snapshotLen int32         = 1024
	promiscuous bool          = false
	timeout     time.Duration = 30 * time.Second

	buffer  gopacket.SerializeBuffer
	options gopacket.SerializeOptions = gopacket.SerializeOptions{FixLengths: true, ComputeChecksums: true}
)*/
// END REMOVE 

// SendAddMapping creates and sends a packet with the following payload:
//			-> 1 byte: Control Type
//				0x01: add mapping
// 		   	-> 4 bytes: from IP
//		   	-> 4 bytes: to IP
// 		   	-> 2 bytes: from Port
//         	-> 2 bytes: to Port
func SendAddMapping(srcIP []byte, dstIP []byte, srcPort []byte, dstPort []byte) {
	var payload []byte; 
	payload = append(payload, 0x01);
	payload = append(payload, srcIP...);
	payload = append(payload, dstIP...);
	payload = append(payload, srcPort...);
	payload = append(payload, dstPort...);

	sendContolPacket(payload);
};

// SendAddDestMapping creates and sends a packet with the following payload:
// 			-> 1 byte: Control Type
//				0x03: add mapping
// 			-> 4 bytes: from IP
//		   	-> 4 bytes: to IP
// 		   	-> 2 bytes: from Port
//         	-> 2 bytes: to Port
func SendAddDestMapping(srcIP []byte, dstIP []byte, srcPort []byte, dstPort []byte) {
	var payload []byte; 
	payload  = append(payload, 0x03);
	payload = append(payload, srcIP...);
	payload = append(payload, dstIP...);
	payload = append(payload, srcPort...);
	payload = append(payload, dstPort...);

	sendContolPacket(payload);
};

// SendListMappings creates and sends a packet with the following payload:
//			-> 1 byte: Control Type
//				0x02: print mappings
func SendListMappings() {
	var payload []byte; 
	payload = append(payload, 0x02);
	sendContolPacket(payload);
};

func createControlPacket(payload []byte) []byte {
	// Start: Not FPGA Friendly 
	// var ipLayer layers.IPv4; 		// does not like "layers.IPv4" 
	// var udpLayer layers.UDP; 		// or "layers.UDP"
	// End 

	// Start: Replacement that runs through parser 
	var ipLayer *IPv4; 
	var udpLayer *UDP; 
	// End

	ipLayer.SrcIP = net.IP(nat.Configs.LAN.IP[:]); 
	ipLayer.DstIP = net.IP(nat.Configs.Ctrl.IP[:]); 
	ipLayer.Version = 4;
	ipLayer.Version = 10;
	ipLayer.Protocol = layers.IPProtocolUDP; 

	udpLayer.SrcPort = layers.UDPPort(80); 
	udpLayer.DstPort = layers.UDPPort(nat.Configs.Ctrl.Port); 
	udpLayer.SetNetworkLayerForChecksum(&ipLayer);

	// Start: Runs through parser 
	var buffer SerializeBuffer; // should be "gopacket.SerializeBuffer" but that does not parse  
	// End
	buffer = gopacket.NewSerializeBuffer(); 	// buffer is a global variables but global variables do not seem parseable 
	gopacket.SerializeLayers(buffer, options,
		&ipLayer,
		&udpLayer,
		gopacket.Payload(payload),
	);
	outgoingPacket := buffer.Bytes();
	return outgoingPacket;
};

func sendContolPacket(payload []byte) {
	handle, err := pcap.OpenLive(nat.Configs.LAN.Name, snapshotLen, promiscuous, timeout);
	if err != nil {
		log.Fatal(err);
	};
	// defer handle.Close();				// defer statement not supported in parse

	outgoingPacket := createControlPacket(payload);
	err = handle.WritePacketData(outgoingPacket);
	if err != nil {
		log.Fatal(err);
	};
};