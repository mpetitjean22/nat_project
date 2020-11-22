package main

import (
	"fmt"
	"io"
	"log"
	"nat_project/pkg/nat"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
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

func setupRawSocketForWANOutbound() int {
	/*
		ifceName := nat.Configs.WAN.Name

		fmt.Println("about to call socket")
		// int(binary.BigEndian.Uint16([]byte{syscall.ETH_P_ALL, 0})) is a hack for htons(ETH_P_ALL)
		fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(binary.BigEndian.Uint16([]byte{syscall.ETH_P_ALL, 0})))
		if err != nil {
			fmt.Println("awk")
			log.Fatal(err)
		}

		fmt.Println("file descriptor", fd)

		err = syscall.SetsockoptString(fd, syscall.SOL_SOCKET, syscall.SO_BINDTODEVICE, ifceName)
		if err != nil {
			fmt.Println("awk2")
			log.Fatal(err)
		}

		// https://stackoverflow.com/questions/22116873/set-socket-option-is-why-so-important-for-a-socket-ip-hdrincl-in-icmp-request
		// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
		// if err != nil {
		// fmt.Println("awk3")
		// log.Fatal(err)
		// }

		return os.NewFile(uintptr(fd), "wan_outbound_raw_socket")
	*/

	ifceName := nat.Configs.WAN.Name

	fmt.Println("about to call socket")
	// int(binary.BigEndian.Uint16([]byte{syscall.ETH_P_ALL, 0})) is a hack for htons(ETH_P_ALL)
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		fmt.Println("awk")
		log.Fatal(err)
	}

	fmt.Println("file descriptor", fd)

	err = syscall.BindToDevice(fd, ifceName)
	if err != nil {
		fmt.Println("awk2")
		log.Fatal(err)
	}

	// https://stackoverflow.com/questions/22116873/set-socket-option-is-why-so-important-for-a-socket-ip-hdrincl-in-icmp-request
	// err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	// if err != nil {
	// fmt.Println("awk3")
	// log.Fatal(err)
	// }

	return fd
	// return os.NewFile(uintptr(fd), "wan_outbound_raw_socket")
}

func main() {
	// BEGIN PROFILING STUFF
	if true {
		f, err := os.Create("profile.pprof")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT)

	go func() {
		sig := <-sigs
		fmt.Println("received signal:")
		fmt.Println(sig)
		done <- true
	}()
	// END PROFILING STUFF

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

	// go listenLAN(readTunIfce, silentMode, staticMode)
	// listenWAN(writeTunIfce, silentMode)

	// BEGIN PROFILING STUFF
	fmt.Println("awaiting signal")
	go listenLAN(readTunIfce, silentMode, staticMode)
	go listenWAN(writeTunIfce, silentMode)
	<-done
	fmt.Println("exiting")
	// END PROFILING STUFF
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
