/* This file implements the structure of the nat table
 * and also has functions that operate on the table.
 */

package nat;

import ("encoding/binary"); 
import ("errors"); 
import ("fmt"); 
import ("sync"); 

// NAT interface which defined some operations that we can perform on
// a particular NAT table

// Start: Not support by FPGA parser 
/* type NAT interface {
	AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte)
	AddDynamicMapping(srcIP [4]byte, srcPort [2]byte, inboundNat NAT)
	GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error)
	HasMapping(srcIP [4]byte, srcPort [2]byte) bool
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

*/ 
// End

// AddDynamicMapping is used to add a mapping from a particular source IP/port
// Example:
// 		SOURCE: 			10.0.0.1 	(port #n) -> 	10.0.2.15 	(port #m)
// 		DESTINATION: 		10.0.2.15 	(port #m) -> 	10.0.0.1 	(port #n)
// port #m could be randomly assigned as an improvement, but for now #m = #n.

// Start: Not supported by FPGA 
// func (nat *Table) AddDynamicMapping(srcIP [4]byte, srcPort [2]byte, inboundNat NAT) {
// End 

// Start: WIll parse but not functional 
func AddDynamicMapping(nat *Table, srcIP [4]byte, srcPort [2]byte, inboundNat NAT) {
// End 
	nat.AddMapping(srcIP, srcPort, Configs.WAN.IP, srcPort);
	inboundNat.AddMapping(Configs.WAN.IP, srcPort, srcIP, srcPort);
};

// AddMapping simply adds a mapping to the table from (srcIP, srcPort) to (dstIP, dstPort)

// Start: Not supported by FPGA 
//func (nat *Table) AddMapping(srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte) {
// End 

// Start: will parse but not functional 
func AddMapping(nat *Table, srcIP [4]byte, srcPort [2]byte, dstIP [4]byte, dstPort [2]byte) {
// End 
	var key, mapping IPv4Address;
	key.ipAdress = srcIP; 
	key.port = srcPort;


	nat.rwMu.Lock(); 
	
	// Start: Not supported by parser 
	// defer nat.rwMu.Unlock()
	// End 

	if nat.natTable == nil {
		nat.natTable = make(map[IPv4Address]IPv4Address);
	};

	mapping.ipAdress = dstIP; 
	mapping.port = dstPort; 
	nat.natTable[key] = mapping; 
};


var ErrNotFound = errors.New("Not Found");

// GetMapping returns the mapping of ip address and port if found, otherwise returns
// an error with not found. A port number of 0 is considered to be a wild card in which
// any port could match to.

// Start: Not FPGA Friendly 
// func (nat *Table) GetMapping(srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error) {
// End

// Start: FPGA Friendly but not functional 
func GetMapping(nat *Table, srcIP [4]byte, srcPort [2]byte) ([4]byte, [2]byte, error) {
// End 
	var key IPv4Address
	var value IPv4Address
	var ok bool

	key.ipAdress = srcIP; 
	key.port = srcPort; 

	nat.rwMu.RLock();

	// Start: Not supported by FPGA 
	// defer nat.rwMu.RUnlock()
	// End 

	value, ok = nat.natTable[key];
	if !ok {
		// check if a wildcard exists
		key.port = [2]byte{0x00, 0x00}
		value, ok = nat.natTable[key];
		if !ok {
			return value.ipAdress, value.port, ErrNotFound;
		};
	};
	return value.ipAdress, value.port, nil;
};

// HasMapping returns the whether a mapping exists for a given srcIP and srcPort
// pair.

// Start: Not FPGA Friendly
// func (nat *Table) HasMapping(srcIP [4]byte, srcPort [2]byte) bool {
// End 

// Start: FPGA friendlt but not functional 
func  HasMapping(nat *Table, srcIP [4]byte, srcPort [2]byte) bool {
// End 

	var key IPv4Address; 

	key.ipAdress = srcIP; 
	key.port = srcPort; 

	nat.rwMu.RLock();
	
	// Start: Not FPGA Friendly 
	// defer nat.rwMu.RUnlock()
	// End 

	_, ok := nat.natTable[key]
	key.port = [2]byte{0, 0}
	_, okWild := nat.natTable[key]

	return ok || okWild
};

// PrettyPrintTable prints the current nat table in a readable format

// Start: Not FPGA Friendly 
// func (nat *Table) PrettyPrintTable() {
// End 

// Start: FPGA Friendly but not functional 
func PrettyPrintTable(nat *Table) {
// End 

	nat.rwMu.RLock();

	// Start: Not supported by FPGA
	// defer nat.rwMu.RUnlock()
	// End 

	fmt.Println("--------------------------");
	for key, value := range nat.natTable {
		prettyPrintAdress(key);
		fmt.Printf(" to ");
		prettyPrintAdress(value);
		fmt.Printf("\n");
	};
	fmt.Println("--------------------------");
};

func prettyPrintAdress(address IPv4Address) {
	fmt.Printf("%d.%d.%d.%d:%d", address.ipAdress[0], address.ipAdress[1], address.ipAdress[2], address.ipAdress[3], binary.BigEndian.Uint16(address.port[:]));
};
