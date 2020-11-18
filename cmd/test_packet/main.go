package main

import (
	"fmt"
	"log"
	"nat_project/pkg/test_packet"
	"time"

	"github.com/google/gopacket/pcap"
)

var (
	//device string = "enp0s3"
	device      string = "tun2"
	snapshotLen int32  = 1024
	promiscuous bool   = false
	err         error
	timeout     time.Duration = 2 * time.Second
	handle      *pcap.Handle
)

func sendPacket(rawPacket []byte) {
	handle, err = pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	err = handle.WritePacketData(rawPacket)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Sending Test Packet")
}

/*
	Small program which uses gopacket to send out packets.
	Useful for debugging + testing and has no real impact on the NAT
	functionality itself.
*/
func main() {
	fmt.Println("Creating Test Packet")
	payload := test_packet.CreateTestPacket([]byte{0x07})
	fmt.Printf("%#v\n", payload)
	sendPacket(payload)
}
