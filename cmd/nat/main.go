package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"nat_project/pkg/get_packets"
	"nat_project/pkg/nat"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/gopacket/pcap"
	"github.com/songgao/water"
)

var (
	snapshotLen int32         = 65535
	promiscuous bool          = false
	timeout     time.Duration = 10 * time.Millisecond
	outboundNat nat.NAT
	inboundNat  nat.NAT
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
	// Setup up SIGINT handler channel for graceful shutdown
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT)

	go func() {
		<-signals // block until SIGINT sent on channel
		done <- true
	}()

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
	defer ifce.Close()

	// file descriptors are thread(/goroutine)-safe per POSIX
	// not sure about *os.File so make a seperate *os.File with the same fd
	file, ok := ifce.ReadWriteCloser.(*os.File)
	if !ok {
		log.Fatalf("water.Interface is not backed by a fd to /dev/tun")
	}

	fd := file.Fd()

	// shared channel for WAN reading
	handle, err := pcap.OpenLive(nat.Configs.WAN.Name, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := get_packets.NewPacketSource(handle)

	// print listening details
	fmt.Printf("WAN Packets on %s\n", nat.Configs.WAN.Name)
	fmt.Printf("LAN Packets on %s\n", nat.Configs.LAN.Name)
	fmt.Printf("Silent Mode: %v\n", silentMode)

	numLANGoRoutines := 1
	numWANGoRoutines := 1

	for i := 0; i < numLANGoRoutines; i++ {
		readTunIfce := os.NewFile(fd, "tunIfce")
		defer readTunIfce.Close()

		go listenLAN(readTunIfce, silentMode, staticMode)
	}

	for i := 0; i < numWANGoRoutines; i++ {
		writeTunIfce := os.NewFile(fd, "tunIfce")
		defer writeTunIfce.Close()

		go listenWAN(packetSource, writeTunIfce, silentMode)
	}

	// block until SIGINT
	<-done
}

func printDestMapping(dstIP [4]byte, srcIP [4]byte, dstPort [2]byte, newDstIP [4]byte, newDstPort [2]byte) {
	fmt.Println("Mapping Found!")
	fmt.Printf("    Original Destination:")
	printAddressPort(dstIP, dstPort)
	fmt.Printf("    	 New Destination:")
	printAddressPort(newDstIP, newDstPort)
}

func printSourceMapping(srcIP [4]byte, dstIP [4]byte, srcPort [2]byte, newSrcIP [4]byte, newSrcPort [2]byte) {
	fmt.Println("Mapping Found!")
	fmt.Printf("    Original Source:")
	printAddressPort(srcIP, srcPort)
	fmt.Printf("    	 New Source:")
	printAddressPort(newSrcIP, newSrcPort)
}

func printAddressPort(ip [4]byte, port [2]byte) {
	fmt.Printf("%d.%d.%d.%d:%d \n", ip[0], ip[1], ip[2], ip[3], binary.BigEndian.Uint16(port[:]))
}

func printOptions() {
	fmt.Println("Options for Running NAT:")
	fmt.Println("   -S")
	fmt.Println("      Silent Mode silences printing out packets when mappings are found")
	fmt.Println("   --static-mapping")
	fmt.Println("      Disables dynamic mapping and only allows for mappings to be added with control packets")
}
