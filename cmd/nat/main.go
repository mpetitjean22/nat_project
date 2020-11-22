package main

import (
	"fmt"
	"io"
	"log"
	"nat_project/pkg/nat"
	"os"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/songgao/water"
)

var (
	snapshotLen int32         = 65535
	promiscuous bool          = false
	timeout     time.Duration = 10 * time.Millisecond
	outboundNat *nat.Table
	inboundNat  *nat.Table
)

func sendPacketPCAP(handle *pcap.Handle, rawPacket []byte) {
	err := handle.WritePacketData(rawPacket)
	if err != nil {
		log.Fatal(err)
	}
}

func sendPacketTun(writeTunIfce io.ReadWriteCloser, rawPacket []byte) {
	writeTunIfce.Write(rawPacket)
}

func main() {
	// configures NAT with settings, will fatalf if something goes wrong
	nat.ConfigureNAT()

	argsWithProg := os.Args
	silentMode := false
	staticMode := false

	if len(argsWithProg) > 1 {
		for i := 1; i < len(argsWithProg); i++ {
			if argsWithProg[i] == "-S" {
				silentMode = true
			} else if argsWithProg[i] == "--static-mapping" {
				staticMode = true
			} else {
				fmt.Printf("Error: %v is an invalid option \n", argsWithProg[i])
				printOptions()
				return
			}
		}
	}

	// Setup NAT tables
	outboundNat = &nat.Table{}
	inboundNat = &nat.Table{}

	// Setup TUN
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = nat.Configs.LAN.Name

	ifce, err := water.New(config)
	if err != nil {
		log.Fatal(err)
	}

	// file descriptors are thread(/goroutine)-safe per POSIX
	// not sure about *os.File so make a seperate *os.File with the same fd
	file, ok := ifce.ReadWriteCloser.(*os.File)
	if !ok {
		log.Fatalf("water.Interface is not backed by a fd to /dev/tun")
	}

	fd := file.Fd()

	readTunIfce := os.NewFile(fd, "tunIfce")
	writeTunIfce := ifce

	go listenLAN(readTunIfce, silentMode, staticMode)
	listenWAN(writeTunIfce, silentMode)
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

func printOptions() {
	fmt.Println("Options for Running NAT:")
	fmt.Println("   -S")
	fmt.Println("      Silent Mode silences printing out packets when mappings are found")
	fmt.Println("   --static-mapping")
	fmt.Println("      Disables dynamic mapping and only allows for mappings to be added with control packets")
}
