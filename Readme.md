# NAT Project
Note: go to tunNAT branch to see project that is operating between two interfaces. 

## Example 
Scroll down for more information about how to run! 

Run packets in one terminal window: 
``` sh 
$ packets 
``` 
In another terminal window, we can add a outbound mapping using control, for example: 
``` sh 
$ control 1 10.0.0.123 0 3.3.3.3 80  
``` 
This will create and send out a control packet which will create a mapping from `10.0.0.123:*` (internal) to `3.3.3.3:80` (external). We use a port of 0 to represent a wild card (*) in which any port will match to it. 

We can also create an inbound mapping using control, for example: 
``` sh 
$ control 3 10.0.0.123 0 4.4.4.4 80  
``` 
This will create and send out a control packet which will create a mapping from `10.0.0.123:*` (external) to `4.4.4.4:80` (internal). 

We can verify that it has been added to the map with the following: 
``` sh 
$ control 3
``` 
which will send a control packet asking to list out all of the current mappings of both the inbound and outbound nat tables. The output will look like the following: 

``` sh 
$ packets
Capturing Packets
Outbound
map[{[10 0 0 123] [0 0]}:0xc0001de020]
Inbound
map[{[10 0 0 123] [0 0]}:0xc000094b60]
``` 

We can verify the functionality by sending out a curl command and having wireshark running in the background. For example, running the curl: 
``` sh 
$ curl 1.1.1.1
```
In wireshark we can see the following: 
```
281	8.595815	3.3.3.3	1.1.1.1	HTTP 
282	8.596620	1.1.1.1	4.4.4.4	TCP
283	8.597428	1.1.1.1	4.4.4.4	TCP
285	8.597929	3.3.3.3	1.1.1.1	TCP	
286	8.598429	1.1.1.1	4.4.4.4	HTTP
```
This shows that for outgoing packets going to `1.1.1.1` are having their source rewritten to `3.3.3.3` which was the rule we set for outgoing packets. 

Conversly, we see that packets coming from `1.1.1.1` have their destination rewritten to `4.4.4.4` which was the rule set for incoming packets. 

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
- implement support for IPv6

- improve `GetMacAddress` function name 
    * remove any hard coded mac address values

- test functionality on the VM 

### FPGA Improvements 
- try and remove := (static variables)
- remove byte slices (this might be a little complicated -- have to see)

