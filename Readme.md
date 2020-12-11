# NAT Project Through FPGA Parser

To view the implementation that currently works please refer to master. 

This branch is intended to see how the code runs through the parser with adjustments to work with the current FPGA parser. In some places code is commented out and replaced with versions that will parse through the parser but that are not functional. As a result, this branch does not run as is! 

--- 
### Runs Through Parser and Functional 
These are files that both run through the parser sucessfully and are fully functional. 

* `./pkg/process_packet/checksum.go`
* `./pkg/process_packet/copy_functions.go`
* `./pkg/process_packet/create_packets.go`
* `./pkg/process_packet/process_packets.go`

* `./cmd/control/main.go`

--- 
### Runs through Parser but needs Modifications to be Functional 
These are files that run through the parser sucessfully but some lines have been modified/removed in order to do so. As a result, functionality has been sacraficed, but with some improvements to the FPGA parser and minor modifications to the code they will work sucessfully. 

Look into each file to see what specifically needs to be modified or improved. I put a lot of effort to be specific about what does not work with the parser. 

* `./pkg/control_packet/process_control_packets.go`
* `./pkg/control_packet/send_control_packets.go`
* `./pkg/get_packets/get_packets.go `
* `./cmd/nat/listen.go`
* `./cmd/nat/main.go`
* `./pkg/nat/nat.go`

## Improvements that are Needed from Parser 

* Supporting argument types that have a `.` in it. 
  * For example, in `./cmd/nat/listen.go` on line 14, `func listenWAN(packetSource *get_packets.PacketSource, writeTunIfce io.ReadWriteCloser, silentMode bool) {` does not parse, but when `get_packets.PacketSource` is modified to just be `PacketSource` it does parse. This is a problem consistantly throughout many files (all instances have comments in their files). 


* Defer statements. I use there frequently and are very handy. 
* Blocks of global variables.
  * For example, on lines 18-22 in `./cmd/nat/main.go` I have a var block which covers multiple variables similar to `import`. When I seperated them out, it would not get through the parser either.

* Passing strings as arguments of functions 
  * On line 155 in `./cmd/nat/maing.go` I have `readTunIfce := os.NewFile(fd, "tunIfce");` which gets the error `primitive type failed for node 1371` when I attempt to run it through the parser. Commenting it out has eliminated the error.

* Creating structs did not seem to parse correctly 
  * For example, see `./pkg/get_packets/get_packets.go` on lines 18-20 

* Creating and having interfaces 
  * For an interface, the function are written as `// func (nat *Table) PrettyPrintTable() {` but this was not parsed properly, instead they were changed to `func PrettyPrintTable(nat *Table) {` which went through the parser but is not fully functional (would have to refactor other parts of the codebase). But I think supporting interfaces is a pretty important part of writing code in Go. 

### Not Intended to Run on FPGA 
These are files which are not intended to run on the FPGA but were instead built to make an easier developer experience on CPU. 

* `./pkg/process_packet/process_packet_test.go`
* `./pkg/control_packet/control_packet_test.go`
* `./pkg/nat/config.go`
* `./pkg/nat/nat_test.go`