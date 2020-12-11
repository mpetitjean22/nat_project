/* Repurposed code from gopacket in order to
 * read from the pcap interface and put into a
 * channel. Was put here in order to reduce the
 * dependency on gopacket as a library (to be used on
 * FPGA down the road)
 */

package get_packets;

import ("io"); 
import ("net"); 
import ("strings"); 
import ("syscall"); 

import ("github.com/google/gopacket/pcap"); 

// Start: Not FPGA Friendly 
/* type PacketSource struct {
	source *pcap.Handle;
	c      chan []byte;
};*/ 
// End

// Start: Not FPGA Friendly 
// func NewPacketSource(source *pcap.Handle) *PacketSource {
// End 

// Start: Runs through Parser (but not functional)
func NewPacketSource(source *Handle) *PacketSource {
// End 
	var result PacketSource; 
	result.source = source; 
	return &result; 
};

// Start: Not FPGA Friendly 
// func (p *PacketSource) NextPacket() ([]byte, error) {

// Start: Runs through Parser (but not functional)
func NextPacket(p *PacketSource) ([]byte, error) {
// End 
	data, _, err := p.source.ReadPacketData();
	if err != nil {
		return nil, err;
	};
	return data, nil;
};

// Start: Not FPGA Friendly 
//func (p *PacketSource) packetsToChannel() {
// End 

// Start: Runs through Parser (but not functional)
func packetsToChannel(p *PacketSource) {
// End 

	// Start: Not FPGA Friendly 
	// defer close(p.c);
	// End 
	for {
		data, err := p.NextPacket();
		if err == nil {
			p.c <- data;
			continue;
		};

		// Start: Not FPGA Friendly 
		// Immediately retry for temporary network errors
		// if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
		//	continue;
		// };
		// End

		// Immediately retry for EAGAIN
		if err == syscall.EAGAIN {
			continue;
		};

		// Immediately break for known unrecoverable errors
		if err == io.EOF || err == io.ErrUnexpectedEOF ||
			err == io.ErrNoProgress || err == io.ErrClosedPipe || err == io.ErrShortBuffer ||
			err == syscall.EBADF ||
			strings.Contains(err.Error(), "use of closed file") {
			break;
		};

		// Sleep briefly and try again
		// time.Sleep(time.Millisecond * time.Duration(5))
	};
};

// Start: Not FPGA friendly 
// func (p *PacketSource) Packets() chan []byte {
// End 

// Start: Runs through parser (but not functional) 
func Packets(p *PacketSource) chan []byte {
// End
	if p.c == nil {
		p.c = make(chan []byte, 10000);
		go p.packetsToChannel();
	};
	return p.c;
};