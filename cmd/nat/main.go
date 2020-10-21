package main

import (
	"fmt"
	"nat_project/pkg/nat"
)

/* Driver for NAT Stuff, if needed */
func main() {
	nat.AddMapping("1.1.1.1", "80", "2.2.2.2", "90")
	nat.AddMapping("1.1.1.1", "60", "2.2.2.2", "40")
	nat.AddMapping("3.3.3.3", "80", "4.4.4.4", "90")

	mappings := nat.ListMappings()
	fmt.Println(mappings)

	dstIP, dstPort, err := nat.GetMapping("1.1.1.1", "60")
	if err != nil {
		fmt.Println("error")
	}
	fmt.Printf("%s %s \n", dstIP, dstPort)
}
