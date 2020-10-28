# NAT Project
## Example 
Scroll down for more information about how to run! 

Run packets in one terminal window: 
``` sh 
$ packets 
``` 
In another terminal window, add a mapping using control, for example: 
``` sh 
$ control 1 10.0.0.252 8009 2.2.2.2 80  
``` 
This will create and send out a control packet which will create a mapping from 10.0.0.252:8009 (internal) to (2.2.2.2:80) external. 

We can verify that it has been added to the map with the following: 
``` sh 
$ control 3
``` 
which will send a control packet asking to list out all of the current mappings, and will result in the following: 

``` sh 
$ packets
Capturing Packets
map[10.0.0.252/8009:2.2.2.2/80]
``` 

On my computer, this mapping is used for some programs internally so we do not need to do anything to see packets with 10.0.0.252/8009 going through and being rewritten. We can see what is happening through our program: 


``` sh 
Mapping Found!
    Original Source: 10.0.0.252:8009
    	 New Source: 2.2.2.2:80
        Destination: 10.0.0.123
```
We can verify that a packet with the new source IP and source port is being written using wireshark: 

``` 
Frame 3450: 176 bytes on wire (1408 bits), 176 bytes captured (1408 bits) on interface en0, id 0

Internet Protocol Version 4, Src: 2.2.2.2, Dst: 10.0.0.123

Transmission Control Protocol, Src Port: 80, Dst Port: 58453, Seq: 1, Ack: 1, Len: 110
``` 

--- 
## How to Run 
### Creating Control Packets 
This will create control packets to send. For now the only control packet 
that this will send will be an `Add Mapping` control packet. However, 
the program is set up to handle other kinds of control packets in the future
with little additional work needed. 
```sh
$ make control
$ control $source ip$, $source port$, $destination ip$, $destination port$      
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
- implement the dual nat tables, one for Internal->External and one for External->Internal
- implement mutex locks on the NAT mapping so that we do not run into any weird situations 
- implement support for IPv6! 
