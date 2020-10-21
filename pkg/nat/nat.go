package nat

import (
	"fmt"
	"strings"
)

var nat_table map[string]string

func AddMapping(srcIP string, srcPort string, dstIP string, dstPort string) {
	if nat_table == nil {
		nat_table = make(map[string]string)
	}
	key := srcIP + "/" + srcPort
	_, ok := nat_table[key]
	if !ok {
		nat_table[key] = ""
	}
	nat_table[key] = dstIP + "/" + dstPort
}

func ListMappings() map[string]string {
	return nat_table
}

func GetMapping(srcIP string, srcPort string) (string, string, error) {
	key := srcIP + "/" + srcPort
	value, ok := nat_table[key]
	if !ok {
		return "", "", fmt.Errorf("Not Found")
	}

	valueArr := strings.Split(value, "/")
	return valueArr[0], valueArr[1], nil
}
