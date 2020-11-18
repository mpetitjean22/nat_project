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
	snapshotLen int32         = 2048
	promiscuous bool          = false
	timeout     time.Duration = 10 * time.Millisecond
	outboundNat *nat.NAT_Table
	inboundNat  *nat.NAT_Table

	tunIfce    *water.Interface
	tunIfceMtx sync.Mutex
	wg         sync.WaitGroup
)

func sendPacketPCAP(handle *pcap.Handle, rawPacket []byte) {
	err := handle.WritePacketData(rawPacket)
	if err != nil {
		log.Fatal(err)
	}
}

func sendPacketTun(rawPacket []byte) {
	// fmt.Println("About to lock in write")
	tunIfceMtx.Lock()
	tunIfce.Write(rawPacket)
	// fmt.Println("About to unlock in write")
	tunIfceMtx.Unlock()
}

func listenWAN(silentMode bool) {
	handle, err := pcap.OpenLive("enp0s3", snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

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

			newIP, newPort, err := inboundNat.GetMapping(dstIP, dstPort)
			if err == nil {
				if bytes.Equal(srcIP[:], []byte{10, 0, 2, 15}) { // TEMP CODE
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
	handle, err := pcap.OpenLive("enp0s3", snapshotLen, promiscuous, timeout) // used for writing
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
			control_packet.ProcessControlPacket(packetData, outboundNat, inboundNat)
		} else {
			newIP, newPort, err := outboundNat.GetMapping(srcIP, srcPort)
			if err == nil {

				if !silentMode || dstIP == [4]byte{1, 2, 3, 4} { // TEMP CODE
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

func main() {
	argsWithProg := os.Args
	silentMode := false
	if len(argsWithProg) == 2 {
		if argsWithProg[1] == "-S" {
			silentMode = true
		}
	}

	outboundNat = &nat.NAT_Table{}
	inboundNat = &nat.NAT_Table{}

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
	go listenLAN(silentMode)
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
