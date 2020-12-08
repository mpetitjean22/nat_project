package nat

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Configurations struct {
	LAN  Interface     `yaml:"LAN-Interface"`
	WAN  Interface     `yaml:"WAN-Interface"`
	Ctrl ControlPacket `yaml:"Control-Packet"`
	NAT  NATSettings   `yaml:"NAT"`
}

type Interface struct {
	Name string  `yaml:"Name"`
	IP   [4]byte `yaml:"IP"`
	Src  [6]byte `yaml:"Src-MAC"`
	Dst  [6]byte `yaml:"Dst-MAC"`
}

type ControlPacket struct {
	IP   [4]byte `yaml:"IP"`
	Port uint16  `yaml:"Port"`
}

type NATSettings struct {
	WANRoutines int `yaml:"WAN-Routines"`
	LANRoutines int `yaml:"LAN-Routines"`
}

var Configs Configurations

func ConfigureNAT() {
	Configs = Configurations{}

	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatalf("Could not find Config YAML")
	}

	err = yaml.Unmarshal(yamlFile, &Configs)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if Configs.LAN.Name == "" || Configs.WAN.Name == "" {
		log.Fatalf("Interface Name cannot be empty string")
	}

	// Debugging Output
	/* fmt.Println(Configs.LAN.Name)
	fmt.Println(Configs.LAN.IP)

	fmt.Println(Configs.WAN.Name)
	fmt.Println(Configs.WAN.IP)
	fmt.Printf("%#v\n", Configs.WAN.Src)
	fmt.Printf("%#v\n", Configs.WAN.Dst)

	fmt.Printf("%#v\n", Configs.Ctrl.IP)
	fmt.Printf("%#v\n", Configs.Ctrl.Port)

	fmt.Printf("%d \n", Configs.NAT.WANRoutines)
	fmt.Printf("%d \n", Configs.NAT.LANRoutines) */
}
