package main

import (
	"fmt"
	"log"
	"nat_project/pkg/control_packet"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
	"nat_project/pkg/process_packet"
	"net"
	"strconv"
	"time"

	"github.com/google/gopacket/pcap"
)

var (
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 2 * time.Second
	handle       *pcap.Handle
)

func main() {
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Capturing Packets")

	packet_source := get_packets.NewPacketSource(handle)

	for packet_data := range packet_source.Packets() {
		ethProtocol, err := process_packet.GetEthProtocol(packet_data)
		if err != nil {
			//fmt.Println(err)
			continue
		}

		if ethProtocol == 0x0800 || ethProtocol == 0x86DD {

			srcIP, dstIP, err := process_packet.GetSrcDstIP(packet_data[14:])
			if err != nil {
				//fmt.Println(err)
				continue
			}

			srcPort, dstPort, err := process_packet.GetSrcDstPort(packet_data[14:])
			if err != nil {
				//fmt.Println(err)
				continue
			}

			controlIP := net.ParseIP("8.8.8.8")
			controlPort := uint16(80)
			if dstIP.Equal(controlIP) && dstPort == controlPort {
				control_packet.ProcessControlPacket(packet_data)
			} else {
				newSrcIP, newSrcPort, err := nat.GetMapping(srcIP.String(), strconv.Itoa(int(srcPort)))
				if err == nil {
					fmt.Println("Mapping Found!")
					fmt.Printf("    Original Source: %s:%d\n", srcIP, srcPort)
					fmt.Printf("    	 New Source: %s:%s \n\n", newSrcIP, newSrcPort)
				}
			}
		}
	}
}
