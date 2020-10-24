package nat

import (
	"fmt"
)

type IPAddress struct {
	ipAdress [4]byte
	port     [2]byte
}

var nat_table map[IPAddress]*IPAddress

func AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte) {
	if nat_table == nil {
		nat_table = make(map[IPAddress]*IPAddress)
	}

	key := IPAddress{
		srcIP,
		srcPort,
	}

	_, ok := nat_table[key]
	if !ok {
		nat_table[key] = &IPAddress{}
	}
	mapping, _ := nat_table[key]
	mapping.ipAdress = dstIP
	mapping.port = dstPort
}

func ListMappings() map[IPAddress]*IPAddress {
	return nat_table
}

func GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error) {
	key := IPAddress{
		srcIP,
		srcPort,
	}
	value, ok := nat_table[key]
	if !ok {
		return [4]byte{}, [2]byte{}, fmt.Errorf("Not Found")
	}
	return value.ipAdress, value.port, nil
}
