package nat

import (
	"fmt"
)

type NAT interface {
	AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte)
	ListMappings() map[IPAddress]*IPAddress
	GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error)
}

type NAT_Table struct {
	nat_table map[IPAddress]*IPAddress
}

type IPAddress struct {
	ipAdress [4]byte
	port     [2]byte
}

func (nat *NAT_Table) AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte) {
	var ok bool
	var mapping *IPAddress
	var key IPAddress

	if nat.nat_table == nil {
		nat.nat_table = make(map[IPAddress]*IPAddress)
	}

	key = IPAddress{
		srcIP,
		srcPort,
	}

	_, ok = nat.nat_table[key]
	if !ok {
		nat.nat_table[key] = &IPAddress{}
	}
	mapping, _ = nat.nat_table[key]
	mapping.ipAdress = dstIP
	mapping.port = dstPort
}

func (nat *NAT_Table) ListMappings() map[IPAddress]*IPAddress {
	return nat.nat_table
}

func (nat *NAT_Table) GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error) {
	var key IPAddress
	var value *IPAddress
	var ok bool

	key = IPAddress{
		srcIP,
		srcPort,
	}
	value, ok = nat.nat_table[key]
	if !ok {
		// check if a wildcard exists (useful for testing)
		key = IPAddress{
			srcIP,
			[2]byte{0x00, 0x00},
		}
		value, ok = nat.nat_table[key]
		if !ok {
			return [4]byte{}, [2]byte{}, fmt.Errorf("Not Found")
		}
	}
	return value.ipAdress, value.port, nil
}
