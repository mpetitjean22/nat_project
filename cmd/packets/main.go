package main

import (
	"bytes"
	"fmt"
	"log"
	"nat_project/pkg/control_packet"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
	"nat_project/pkg/process_packet"
	"os"
	"sync"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/songgao/water"
)

var (
	snapshot_len int32         = 2048
	promiscuous  bool          = false
	timeout      time.Duration = 10 * time.Millisecond
	outbound_nat *nat.NAT_Table
	inbound_nat  *nat.NAT_Table

	tunIfce    *water.Interface
	tunIfceMtx sync.Mutex
	wg         sync.WaitGroup
)

func sendPacket(handle *pcap.Handle, rawPacket []byte) {
	//packet := []byte{0x52, 0x54, 0x00, 0x12, 0x35, 0x02, 0x08, 0x00, 0x27, 0xfd, 0x06, 0x32, 0x08, 0x00, 0x45, 0x00, 0x00, 0x3c, 0x04, 0x70, 0x40, 0x00, 0x40, 0x06, 0x28, 0x3c, 0x0a, 0x00, 0x02, 0x0f, 0x01, 0x01, 0x01, 0x01, 0xe0, 0x6a, 0x00, 0x50, 0xc1, 0xa1, 0x83, 0x9b, 0x00, 0x00, 0x00, 0x00, 0xa0, 0x02, 0xfa, 0xf0, 0x0e, 0x3f, 0x00, 0x00, 0x02, 0x04, 0x05, 0xb4, 0x04, 0x02, 0x08, 0x0a, 0x15, 0xbd, 0x50, 0xd1, 0x00, 0x00, 0x00, 0x00, 0x01, 0x03, 0x03, 0x07}

	err := handle.WritePacketData(rawPacket)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%#v\n", rawPacket)
	//fmt.Println("Sending out on enp0s3")
}

func sendPacketTun(rawPacket []byte) {
	// fmt.Println("About to lock in write")
	tunIfceMtx.Lock()
	tunIfce.Write(rawPacket)
	// fmt.Println("About to unlock in write")
	tunIfceMtx.Unlock()
}

func listenWAN(silentMode bool) {
	handle, err := pcap.OpenLive("enp0s3", snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Capturing Packets on enp0s3")
	fmt.Printf("Silent Mode: %v \n", silentMode)

	packetSource := get_packets.NewPacketSource(handle)

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

			newIP, newPort, err := inbound_nat.GetMapping(dstIP, dstPort)
			if err == nil {
				if bytes.Equal(srcIP[:], []byte{10, 0, 2, 15}) {
					if !silentMode || true {
						printDestMapping(dstIP, srcIP, dstPort, newIP, newPort)
					}

					newPacketData, err := process_packet.WriteDestination(packetData, newIP, newPort)
					if err == nil {
						sendPacketTun(newPacketData[14:len(packetData)])
					}
				}
			}
		}
	}
}

func listenLAN(silentMode bool) {
	// Open device
	/*handle, err := pcap.OpenLive("tun2", snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	if err != nil {
		log.Fatal(err)
	}*/

	fmt.Println("Capturing Packets on tun2")
	fmt.Printf("Silent Mode: %v \n", silentMode)

	//packetSource := get_packets.NewPacketSource(handle)

	/*for packetData := range packetSource.Packets() {

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

		if dstIP == control_packet.ControlIP && dstPort == control_packet.ControlPort {
			control_packet.ProcessControlPacket(packetData, outbound_nat, inbound_nat)
		} else {
			newIP, newPort, err := outbound_nat.GetMapping(srcIP, srcPort)
			if err == nil {

				if !silentMode || dstIP == [4]byte{1, 2, 3, 4} {
					printSourceMapping(srcIP, dstIP, srcPort, newIP, newPort)
				}

				newPacketData, err := process_packet.WriteSource(packetData, newIP, newPort)
				if err == nil {
					sendPacket("enp0s3", newPacketData[:len(packetData)+14])
				}
			}
		}
	}*/
}

func listenLAN2(silentMode bool) {
	handle, err := pcap.OpenLive("enp0s3", snapshot_len, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	buffer := make([]byte, 2048)
	fmt.Println("Capturing Packets on tun2")
	fmt.Printf("Silent Mode: %v \n", silentMode)

	for {
		// fmt.Println("About to lock in read")
		tunIfceMtx.Lock()
		n, err := tunIfce.Read(buffer)
		// fmt.Println("About to unlock in read")
		tunIfceMtx.Unlock()
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

		if dstIP == control_packet.ControlIP && dstPort == control_packet.ControlPort {
			control_packet.ProcessControlPacket(packetData, outbound_nat, inbound_nat)
		} else {
			newIP, newPort, err := outbound_nat.GetMapping(srcIP, srcPort)
			if err == nil {

				if !silentMode || dstIP == [4]byte{1, 2, 3, 4} {
					printSourceMapping(srcIP, dstIP, srcPort, newIP, newPort)
				}

				newPacketData, err := process_packet.WriteSource(packetData, newIP, newPort)
				if err == nil {
					sendPacket(handle, newPacketData[:len(packetData)+14])
				}
			}
		}
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

	outbound_nat = &nat.NAT_Table{}
	inbound_nat = &nat.NAT_Table{}

	// Setup TUN
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = "tun2"

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}
	tunIfce = ifce

	wg.Add(2)
	go listenLAN2(silentMode)
	go listenWAN(silentMode)
	wg.Wait()
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
