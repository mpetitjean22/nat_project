# NAT Project

### Table of Contents
- [Demo Instructions](#Demo)
- [Configuations](#Configuations)
- [NAT](#NAT)
    * [Capturing Packets](#Capturing-Packets)
    * [Options](#Running-NAT)
- [Control Packets](#Control-Packets)
    * [Control CLI](#Control-CLI)
- [Tests](#Tests)
- [Remaining Work](#Remaining-Work)
    * [General Improvements](#General-Improvements)
    * [FPGA Improvements](#FPGA-Improvements)


## Demo 

Make and run `nat`. The `-S` option can be used to put the nat in silent mode. 
``` sh 
$ make nat
# if you get "command not found" run "source ./scripts/add-to-path.sh"
# without silent mode: 
$ sudo $(which nat)
# with silent mode: 
$ sudo $(which nat) -S
``` 

Configure the TUN interface: 
``` sh 
$ source scripts/set-tun.sh
``` 

Add IP Routes for google and wikipedia, which will run packets through the NAT. 
``` sh 
$ source scripts/demo.sh
```

Run a curl command to google: 
```sh 
$ sudo curl --verbose --interface tun2 -ipv4 https://www.google.com
``` 

Open lynx and browse google and wikipedia: 
```sh 
$ lynx google.com
``` 

Optionally open `tcpdump` to view packets passing through the two interfaces! 

--- 
## Configuations 
In order to make the NAT and Control programs easier to use, the `config.yaml` file allows the user to easily change and specify configuations. 

``` 
LAN-Interface:          # specifies tun interface
  Name: tun2            # name of tun interface that is created
  IP: [10,0,0,1]        # ip address of the tunu interface 
  
WAN-Interface:          # specifies eth interface 
  Name: enp0s3
  IP: [10,0,2,15]
  Src-MAC: [0x52,0x54,0x00,0x12,0x35,0x02]
  Dst-MAC: [0x08,0x00,0x27,0xfd,0x06,0x32]

Control-Packet: 
  IP: [10,0,0,2]        # dest IP of the control packets 
  Port: 80              # dest Port of the control packets 

NAT: 
  WAN-Routines: 1       # number of routines (still in testing)
  LAN-Routines: 1
```
It is important to note that these values are not validated for correctness so it is up to the user to provide configurations that are correct. Additionally, `tun2` is the tun name used in the scripts, so if you change this value be sure to change the scripts or do not use them. 

--- 
## NAT 
### Capturing Packets 
The NAT listens on two interfaces. The TUN (tun2) interface is considered the LAN side of the network while the eth interface (enp0s3) is considered the WAN side of the network. 

By default, the NAT will use dynamic mapping in order to create mappings from a specific LAN IP and Port to a WAN IP and Port. The NAT also takes control packets which allows the user to add static mappings for both the WAN and LAN side of the NAT. 

### Running NAT 
First you must make the nat and make sure that you have root privileges. 
``` sh 
$ make nat
# if you get "command not found" run "source ./scripts/add-to-path.sh"
$ sudo $(which nat)
```

You can run the NAT with a couple options which are specified as: 
``` sh
Options for Running NAT:
   -S
      Silent Mode silences printing out packets when mappings are found
   --static-mapping
      Disables dynamic mapping and only allows for mappings to be added with control packets
```

The NAT will run continuously until it is closed with `^C`. 

---
## Control Packets 

Control is a command line tool which allows us to create and send fully formed control packets to the NAT program. The control packets are IPv4 on top of UDP, with a special payload format. 

IPv4 Header:
* Source IP: 8.8.8.8

UDP Header: 
* Source Port: 80

Payload: 
* Control Type (1 Byte)
    * 0x01 = Add Mapping (LAN -> WAN) 
    * 0x02 = List Mappings 
    * 0x03 = Add Mapping (WAN -> LAN) 
* From IP (4 Bytes) 
* To IP (4 Bytes) 
* From Port (2 Bytes)
    * A source port of 0 will be considered as a wildcard and will only match the source IP 
* To Port (2 Bytes) 

### Control CLI 

The program provides a Command Line Interface for creating and sending control packets with the proper format as described above.
```sh
$ make control
 # if you get "command not found" run "source ./scripts/add-to-path.sh"
$ sudo $(which control)
Error: Invalid Number of Arguments
Looking for: (control type)
   control types:
       2 -> list current mappings
Looking for: (control type) (fromIP) (fromPort) (toIP) (toPort)
   control types:
       1 -> add outbound mapping
       3 -> add inbound mapping
```

Sending a control packet to list mappings would look like: 
``` sh 
$ sudo $(which control) 2
```

Sending a packet to create a LAN to WAN mapping would look like: 
``` sh 
$ sudo $(which control) 1 1.1.1.1 80 2.2.2.2 80 
```
--- 
## Tests
There are test files implemented in order to verify the functionality of every part of the project. They are located in each of the subdirectories in `pkg`, and end with `_test.go`. 

---
## Remaining Work
### General Improvements
- write up on sakai (in progress) 

### FPGA Improvements 
- make a branch + run through the parser and put the results 

