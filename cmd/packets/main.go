package main

import (
	"fmt"
	"log"
	"nat_project/pkg/control_packet"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
	"nat_project/pkg/process_packet"
	"os"
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

func sendPacket(rawPacket []byte) {
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	err = handle.WritePacketData(rawPacket)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	argsWithProg := os.Args
	silentMode := false
	if len(argsWithProg) == 2 {
		if argsWithProg[1] == "-S" {
			silentMode = true
		}
	}

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
	fmt.Printf("Silent Mode: %v \n", silentMode)

	packet_source := get_packets.NewPacketSource(handle)
	outbound_nat := nat.NAT_Table{}
	inbound_nat := nat.NAT_Table{}

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

			controlIP := [4]byte{0x08, 0x08, 0x08, 0x08}
			controlPort := [2]byte{0x00, 0x50}

			if dstIP == controlIP && dstPort == controlPort {
				control_packet.ProcessControlPacket(packet_data, &outbound_nat, &inbound_nat)
			} else {
				if !process_packet.GetMacAddress(packet_data) {
					newSrcIP, newSrcPort, err := outbound_nat.GetMapping(srcIP, srcPort)
					if err == nil {

						if !silentMode {
							printSourceMapping(srcIP, dstIP, srcPort, newSrcIP, newSrcPort)
						}

						newPacketData, err := process_packet.WriteSource(packet_data, newSrcIP, newSrcPort)
						if err == nil {
							sendPacket(newPacketData[:len(packet_data)])
						}
					}
				} else {
					newDstIP, newDstPort, err := inbound_nat.GetMapping(dstIP, dstPort)
					if err == nil {

						if !silentMode {
							printDestMapping(dstIP, srcIP, dstPort, newDstIP, newDstPort)
						}

						newPacketData, err := process_packet.WriteDestination(packet_data, newDstIP, newDstPort)
						if err == nil {
							sendPacket(newPacketData[:len(packet_data)])
						}
					}
				}
			}
		}
	}
}

func printDestMapping(dstIP [4]byte, srcIP [4]byte, dstPort [2]byte, newDstIP [4]byte, newDstPort [2]byte) {
	fmt.Println("Mapping Found!")
	fmt.Printf("    Original Destination: %v:%v\n", dstIP, dstPort)
	fmt.Printf("    	 New Destination: %v:%v \n", newDstIP, newDstPort)
	fmt.Printf("                   Source: %v \n \n", srcIP)
}

func printSourceMapping(srcIP [4]byte, dstIP [4]byte, srcPort [2]byte, newSrcIP [4]byte, newSrcPort [2]byte) {
	fmt.Println("Mapping Found!")
	fmt.Printf("    Original Source: %v:%v\n", srcIP, srcPort)
	fmt.Printf("    	 New Source: %v:%v \n", newSrcIP, newSrcPort)
	fmt.Printf("        Destination: %v \n \n", dstIP)
}
