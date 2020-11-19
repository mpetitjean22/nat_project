# NAT Project
## Demo 

Run packets (in silent mode!!) in one terminal window: 
``` sh 
$ make packets
$ packets -S
``` 

Configure the TUN and VM
``` sh 
$ source scripts/set-tun.sh
``` 

Optionally try and send out a curl over tun2 iterface: 
``` sh 
$ sudo curl --verbose --interface tun2 1.1.1.1
``` 
Running wireshark/tcpdump on tun2 interface will show syn packets going out but not response coming in. 

Add mappings to NAT for demo (this will just send control packets -- this is also specific to Marie's VM fyi)
``` sh 
$ source scripts/add-mapping.sh
```

Now when we curl: 
``` sh 
$ sudo curl --verbose --interface tun2 1.1.1.1
```
We will get the response! We can also try this out by attempting to load the google home page: 

```sh 
$ sudo curl --verbose --interface tun2 https://wwww.google.com
``` 


--- 
## How to Run 
### Creating Control Packets 

```sh
$ make control
$ control
Invalid Number of Arguments
Looking for: (control type)
Looking for: (control type) (sourceIP) (sourcePort) (destinationIP) (destinationPort)
   control types:
       1 -> add outbound mapping
       2 -> list current mappings
       3 -> add inbound mapping
 # if you get "command not found" run "source ./scripts/add-to-path.sh"
```
Control packets currently have the following structure: 

IPv4 Header: 
* Source IP: 8.8.8.8

UDP Header: 
* Source Port: 80

Payload: 
* Control Type (1 Byte)
    * 0x01 = Add Mapping (Internal -> External) 
    * 0x02 = List Mappings 
    * 0x03 = Add Mapping (External -> Internal) 
* Source IP (4 Bytes) 
* Destination IP (4 Bytes) 
* Source Port (2 Bytes)
    * Note: A source port of 0 will be considered as a wildcard and will only match the source IP 
* Destination Port (2 Bytes) 

Example/ This will add a mapping from 1.1.1.1/80 to 2.2.2.2/80 in the NAT table. 
```sh
$ make control
$ control 1 1.1.1.1 80 2.2.2.2 80 
```

### Capturing Packets 
This will capture packets and filter between being a control packet and a packet whose IP/Port need to rewritten. Right now, we only consider packets with IPv4/IPv6 and UDP/TCP. However, we do not currently get the ports for the IPv6 header. 

For now, packets are not being rewritten and sent out. But they are being detected and a print statement is made when we
detect a packet that has a mapping and what the mapping is. 

```sh
$ make packets
$ packets
 # if you get "command not found" run "source ./scripts/add-to-path.sh"
```

### Testing NAT Program 
In order to better test and develope the actual NAT part of the project we can use `nat`. 

```sh
$ make nat
$ nat
 # if you get "command not found" run "source ./scripts/add-to-path.sh"
```

This is just as a sanity check and for checking the proper operation. 

---
## Tests
In addition, there are test files implemented in order to test the functionality of every part of the project. The are located in each of the subdirectories in `pkg`, and end with `_test.go`. 

---

## Left Todo
### General Improvements
- implement mutex locks on the NAT mapping so that we do not run into any weird situations 
- NAT is much faster now, but still a little bit slow loading google home page...maybe stop using pcap for injecting packets??
- dynamic mappings!!! 


### FPGA Improvements 
- try and remove := (static variables)
- remove byte slices (this might be a little complicated -- have to see)

