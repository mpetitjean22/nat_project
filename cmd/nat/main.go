package main;

import ("encoding/binary"); 
import ("fmt"); 
import ("io"); 
import ("log"); 
import ("nat_project/pkg/get_packets"); 
import ("nat_project/pkg/nat"); 
import ("os"); 
import ("os/signal"); 
import ("syscall"); 
import ("time"); 

import ("github.com/google/gopacket/pcap"); 
import ("github.com/songgao/water"); 

// Start: Not FPGA Friendly 
//var (snapshotLen int32 = 65535);
//var (promiscuous bool = false); 
//var (timeout     time.Duration = 10 * time.Millisecond); 
//var (outboundNat nat.NAT); 
//var (inboundNat  nat.NAT);
// END 

// Start: Not FPGA Friendly 
// func sendPacketPCAP(handle *pcap.Handle, rawPacket []byte) {
// End 

// Start: FPGA Friendly Version (not functional) 
func sendPacketPCAP(handle *Handle, rawPacket []byte) {
// End 
	err := handle.WritePacketData(rawPacket);
	if err != nil {
		log.Fatal(err);
	};
};

// Start: Not FPGA Friendly 
// func sendPacketTun(writeTunIfce io.ReadWriteCloser, rawPacket []byte) {
// End 

// Start: FPGA Friendly (but not functional) 
func sendPacketTun(writeTunIfce ReadWriteCloser, rawPacket []byte) {
// End 
	writeTunIfce.Write(rawPacket);
};

func main() {
	// Setup up SIGINT handler channel for graceful shutdown
	signals := make(chan Signal, 1);
	done := make(chan bool, 1);

	signal.Notify(signals, syscall.SIGINT);

	go func() {
		<-signals; // block until SIGINT sent on channel
		done <- true;
	}();

	// configures NAT with settings, will fatalf if something goes wrong
	nat.ConfigureNAT();

	argsWithProg := os.Args;
	silentMode := false;
	staticMode := false;

	if len(argsWithProg) > 1 {
		for i := 1; i < len(argsWithProg); i++ {
			if argsWithProg[i] == "-S" {
				silentMode = true;
			} else if argsWithProg[i] == "--static-mapping" {
				staticMode = true;
			} else {
				fmt.Printf("Error: %v is an invalid option \n", argsWithProg[i]);
				printOptions();
				return;
			};
		};
	};

	// Setup NAT tables
	// Start: Not FPGA Friendly 
	//outboundNat = &nat.Table{};
	//inboundNat = &nat.Table{};
	// End 

	// Start: More FPGA Friendly, but some issues with nat.Table
	// var outboundNat nat.Table; 
	// var inboundNat nat.Table; 
	// End 

	// Start: Will parse but not functional 
	var outboundNat Table; 
	var inboundNat Table; 
	// End


	// Setup TUN
	// Start: Not FPGA Friendly
	// config := water.Config{
	// 	DeviceType: water.TUN,
	// };
	// End 

	// Start: More FPGA Friendly, but some issues with "water.TUN"
	// var config water.Config 
	// End 

	// Start: FPGA friendly, but not functional
	var config Config; 
	// End 

	config.Name = nat.Configs.LAN.Name;

	ifce, err := water.New(config);
	if err != nil {
		log.Fatal(err);
	};

	// Start: Not FPGA Friendly 
	// defer ifce.Close();
	// End 
 
	// file descriptors are thread(/goroutine)-safe per POSIX
	// not sure about *os.File so make a seperate *os.File with the same fd
	
	// Start: Not FPA friendly with the format of "something.something"
	// file, ok := ifce.ReadWriteCloser.(*os.File);
	// End 
	if !ok {
		log.Fatalf("water.Interface is not backed by a fd to /dev/tun");
	};

	fd := file.Fd();

	// shared channel for WAN reading
	handle, err := pcap.OpenLive(nat.Configs.WAN.Name, snapshotLen, promiscuous, timeout);
	if err != nil {
		log.Fatal(err);
	};
	//defer handle.Close();

	packetSource := get_packets.NewPacketSource(handle);

	// print listening details
	fmt.Printf("WAN Packets on %s\n", nat.Configs.WAN.Name);
	fmt.Printf("LAN Packets on %s\n", nat.Configs.LAN.Name);
	fmt.Printf("Silent Mode: %v\n", silentMode);

	numLANGoRoutines := nat.Configs.NAT.LANRoutines;
	numWANGoRoutines := nat.Configs.NAT.WANRoutines;

	for i := 0; i < numLANGoRoutines; i++ {
		// Start: Not FPGA Friendly 
		// readTunIfce := os.NewFile(fd, "tunIfce");		// Error produced: primitive type failed for node 1371
		// End

		// Start: Not FPGA friendly 
		//defer readTunIfce.Close();
		// END 

		go listenLAN(readTunIfce, silentMode, staticMode);
	};

	for i := 0; i < numWANGoRoutines; i++ {
		// Start: Not FPGA Friendly 
		// writeTunIfce := os.NewFile(fd, "tunIfce");	//Error produced: primitive type failed for node 1455
		// End 

		// Start: Not FPGA friendly 
		// defer writeTunIfce.Close();
		// End 

		go listenWAN(packetSource, writeTunIfce, silentMode);
	};

	// block until SIGINT
	<-done; 
};

func printDestMapping(dstIP [4]byte, srcIP [4]byte, dstPort [2]byte, newDstIP [4]byte, newDstPort [2]byte) {
	fmt.Println("Mapping Found!");
	fmt.Printf("    Original Destination:");
	printAddressPort(dstIP, dstPort);
	fmt.Printf("    	 New Destination:");
	printAddressPort(newDstIP, newDstPort);
};

func printSourceMapping(srcIP [4]byte, dstIP [4]byte, srcPort [2]byte, newSrcIP [4]byte, newSrcPort [2]byte) {
	fmt.Println("Mapping Found!");
	fmt.Printf("    Original Source:");
	printAddressPort(srcIP, srcPort);
	fmt.Printf("    	 New Source:");
	printAddressPort(newSrcIP, newSrcPort);
};

func printAddressPort(ip [4]byte, port [2]byte) {
	fmt.Printf("%d.%d.%d.%d:%d \n", ip[0], ip[1], ip[2], ip[3], binary.BigEndian.Uint16(port[:]));
};

func printOptions() {
	fmt.Println("Options for Running NAT:");
	fmt.Println("   -S");
	fmt.Println("      Silent Mode silences printing out packets when mappings are found");
	fmt.Println("   --static-mapping");
	fmt.Println("      Disables dynamic mapping and only allows for mappings to be added with control packets");
};