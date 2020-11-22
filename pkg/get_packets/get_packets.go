/* Repurposed code from gopacket in order to
 * read from the pcap interface and put into a
 * channel. Was put here in order to reduce the
 * dependency on gopacket as a library (to be used on
 * FPGA down the road)
 */

package get_packets

import (
	"io"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/google/gopacket/pcap"
)

type PacketSource struct {
	source *pcap.Handle
	c      chan []byte
}

func NewPacketSource(source *pcap.Handle) *PacketSource {
	return &PacketSource{
		source: source,
	}
}

func (p *PacketSource) NextPacket() ([]byte, error) {
	data, _, err := p.source.ReadPacketData()
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (p *PacketSource) packetsToChannel() {
	defer close(p.c)
	for {
		data, err := p.NextPacket()
		if err == nil {
			p.c <- data
			continue
		}

		// Immediately retry for temporary network errors
		if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
			continue
		}

		// Immediately retry for EAGAIN
		if err == syscall.EAGAIN {
			continue
		}

		// Immediately break for known unrecoverable errors
		if err == io.EOF || err == io.ErrUnexpectedEOF ||
			err == io.ErrNoProgress || err == io.ErrClosedPipe || err == io.ErrShortBuffer ||
			err == syscall.EBADF ||
			strings.Contains(err.Error(), "use of closed file") {
			break
		}

		// Sleep briefly and try again
		time.Sleep(time.Millisecond * time.Duration(5))
	}
}

func (p *PacketSource) Packets() chan []byte {
	if p.c == nil {
		p.c = make(chan []byte, 1000)
		go p.packetsToChannel()
	}
	return p.c
}
