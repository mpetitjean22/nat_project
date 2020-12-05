/* This file implements the structure of the nat table
 * and also has functions that operate on the table.
 */

package nat

import (
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
)

// NAT interface which defined some operations that we can perform on
// a particular NAT table
type NAT interface {
	AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte)
	AddDynamicMapping(srcIP [4]byte, srcPort [2]byte, inboundNat NAT)
	GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error)
	PrettyPrintTable()
}

// Table is the struct which holds the NAT Table data structure.
type Table struct {
	natTable map[IPv4Address]IPv4Address
	rwMu     sync.RWMutex
}

// IPv4Address is a struct to hold fixed length array for IP Address and Port
// designed to be compatible with FPGA which does not support slices
type IPv4Address struct {
	ipAdress [4]byte
	port     [2]byte
}

// AddDynamicMapping is used to add a mapping from a particular source IP/port
// Example:
// 		SOURCE: 			10.0.0.1 	(port #n) -> 	10.0.2.15 	(port #m)
// 		DESTINATION: 		10.0.2.15 	(port #m) -> 	10.0.0.1 	(port #n)
// port #m could be randomly assigned as an improvement, but for now #m = #n.
func (nat *Table) AddDynamicMapping(srcIP [4]byte, srcPort [2]byte, inboundNat NAT) {
	//key := IPv4Address{
	//	srcIP,
	//	srcPort,
	//}

	// nat.rwMu.RLock()
	// _, hasMapping := nat.natTable[key]
	// nat.rwMu.RUnlock()

	// if !hasMapping {
	nat.AddMapping(srcIP, srcPort, Configs.WAN.IP, srcPort)
	inboundNat.AddMapping(Configs.WAN.IP, srcPort, srcIP, srcPort)
	// }
}

// AddMapping simply adds a mapping to the table from (srcIP, srcPort) to (dstIP, dstPort)
func (nat *Table) AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte) {
	//var mapping IPv4Address
	var key IPv4Address

	key = IPv4Address{
		srcIP,
		srcPort,
	}

	nat.rwMu.Lock()
	defer nat.rwMu.Unlock()

	if nat.natTable == nil {
		nat.natTable = make(map[IPv4Address]IPv4Address)
	}

	//mapping, hasMapping := nat.natTable[key]
	/* if !hasMapping {
		mapping = IPv4Address{}
		nat.natTable[key] = mapping
	} else {

	}*/
	nat.natTable[key] = IPv4Address{
		dstIP,
		dstPort,
	}
}

var ErrNotFound = errors.New("Not Found")

// GetMapping returns the mapping of ip address and port if found, otherwise returns
// an error with not found. A port number of 0 is considered to be a wild card in which
// any port could match to.
func (nat *Table) GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error) {
	var key IPv4Address
	var value IPv4Address
	var ok bool

	key = IPv4Address{
		srcIP,
		srcPort,
	}

	nat.rwMu.RLock()
	defer nat.rwMu.RUnlock()

	value, ok = nat.natTable[key]
	if !ok {
		// check if a wildcard exists
		key.port = [2]byte{0x00, 0x00}
		value, ok = nat.natTable[key]
		if !ok {
			return value.ipAdress, value.port, ErrNotFound
		}
	}
	return value.ipAdress, value.port, nil
}

// PrettyPrintTable prints the current nat table in a readable format
func (nat *Table) PrettyPrintTable() {
	nat.rwMu.RLock()
	defer nat.rwMu.RUnlock()

	fmt.Println("--------------------------")
	for key, value := range nat.natTable {
		prettyPrintAdress(key)
		fmt.Printf(" to ")
		prettyPrintAdress(value)
		fmt.Printf("\n")
	}
	fmt.Println("--------------------------")
}

func prettyPrintAdress(address IPv4Address) {
	fmt.Printf("%d.%d.%d.%d:%d", address.ipAdress[0], address.ipAdress[1], address.ipAdress[2], address.ipAdress[3], binary.BigEndian.Uint16(address.port[:]))
}
