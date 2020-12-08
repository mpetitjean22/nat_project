package main

import (
	"encoding/binary"
	"io"
	"log"
	"nat_project/pkg/control_packet"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
	"nat_project/pkg/process_packet"

	"github.com/google/gopacket/pcap"
)

func listenWAN(packetSource *get_packets.PacketSource, writeTunIfce io.ReadWriteCloser, silentMode bool) {
	for packetData := range packetSource.Packets() {

		ethProtocol, err := process_packet.GetEthProtocol(packetData)
		if err != nil {
			//fmt.Println(err)
			continue
		}

		if ethProtocol == 0x0800 || ethProtocol == 0x86DD {

			srcIP, dstIP, err := process_packet.GetSrcDstIP(packetData[14:])
			if err != nil {
				//fmt.Println(err)
				continue
			}

			_, dstPort, err := process_packet.GetSrcDstPort(packetData[14:])
			if err != nil {
				//fmt.Println(err)
				continue
			}

			newIP, newPort, err := inboundNat.GetMapping(dstIP, dstPort)
			if err == nil {
				if !silentMode {
					printDestMapping(dstIP, srcIP, dstPort, newIP, newPort)
				}

				newPacketData, err := process_packet.WriteDestination(packetData, newIP, newPort)
				if err == nil {
					sendPacketTun(writeTunIfce, newPacketData[14:len(packetData)])
				}
			}
		}
	}
}

func listenLAN(readTunIfce io.ReadWriteCloser, silentMode bool, staticMode bool) {
	handle, err := pcap.OpenLive(nat.Configs.WAN.Name, snapshotLen, promiscuous, timeout) // used for writing
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	buffer := make([]byte, 65535)

	for {
		n, err := readTunIfce.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		packetData := buffer[:n]

		srcIP, dstIP, err := process_packet.GetSrcDstIP(packetData)
		if err != nil {
			//fmt.Println(err)
			continue
		}

		srcPort, dstPort, err := process_packet.GetSrcDstPort(packetData)
		if err != nil {
			//fmt.Println(err)
			continue
		}

		if dstIP == nat.Configs.Ctrl.IP && binary.BigEndian.Uint16(dstPort[:]) == nat.Configs.Ctrl.Port {
			control_packet.ProcessControlPacket(packetData, outboundNat, inboundNat)
		} else {
			mappingExists := outboundNat.HasMapping(srcIP, srcPort)
			if !mappingExists && !staticMode {
				outboundNat.AddDynamicMapping(srcIP, srcPort, inboundNat)
			}

			newIP, newPort, err := outboundNat.GetMapping(srcIP, srcPort)
			if err == nil {
				if !silentMode {
					printSourceMapping(srcIP, dstIP, srcPort, newIP, newPort)
				}

				newPacketData, err := process_packet.WriteSource(packetData, newIP, newPort)
				if err == nil {
					sendPacketPCAP(handle, newPacketData[:len(packetData)+14])
				}
			}
		}
	}
}
