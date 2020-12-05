# httperf Performance Results 
## Process 
I first set up a local HTTP server on my machine while the NAT was running inside of a VM. 

``` sh 
# Start HTTP Server
$ python3 -m http.server 8080 --bind 0.0.0.0
```

In order to test how the NAT responds to having to create many connections, I used httperf to gather information on my local HTTP server. This test will create a total of 8000 connections, and stops once all of the connections are either completed or a failed. A connection is considered to be failed if it any activity on the connection fails to make forward progress after 20 seconds (the timeout). Connections are created at a rate of 100 per second. 

``` sh 
# Run from inside the VM 
httperf --server 10.0.0.123 --port 8080 --verbose --rate 100 --num-conn 8000 --timeout 20
```

When we want the httperf to run through the NAT, then we must set an ip route to route packets through the tun2 interface. 
``` sh 
$ sudo ip route add 10.0.0.123 dev tun2
``` 

## Conclusion 

### Connection Times
``` sh
# No NAT 
Connection rate: 100.0 conn/s (10.0 ms/conn, <=51 concurrent connections)
Connection time [ms]: min 2.6 avg 6.7 max 1152.6 median 4.5 stddev 23.4
Connection time [ms]: connect 1.7

# NAT
Connection rate: 100.0 conn/s (10.0 ms/conn, <=20 concurrent connections)
Connection time [ms]: min 34.8 avg 48.8 max 1197.5 median 46.5 stddev 23.8
Connection time [ms]: connect 14.5
```

We can see that the NAT is able to maintain sending out 100 connections per second similar to without a NAT. However, the number of concurrent connections is about half, 21 concurrect connection vs 51 concurrent connections. This is likely due to the NAT introducing some latency which reduces the number of concurrent connections possible. 

Over the 8000 connections, the NAT introduced an additional 42.1 milliseconds per connection. The distribution of connection rates look similar, with similar stdevs and maximum values. The only increase appears in the minimum connection time and also average/median. 

### Reply Times
``` sh 
# No NAT
Reply rate [replies/s]: min 91.4 avg 100.0 max 108.8 stddev 3.3 (15 samples)
Reply time [ms]: response 4.6 transfer 0.3
Reply size [B]: header 155.0 content 4393.0 footer 0.0 (total 4548.0)
Reply status: 1xx=0 2xx=8000 3xx=0 4xx=0 5xx=0

# NAT
Reply rate [replies/s]: min 99.0 avg 100.0 max 100.8 stddev 0.4 (16 samples)
Reply time [ms]: response 18.3 transfer 15.9
Reply size [B]: header 155.0 content 4393.0 footer 0.0 (total 4548.0)
Reply status: 1xx=0 2xx=8000 3xx=0 4xx=0 5xx=0
```

The reply rate with and without the NAT are very similar and show an average of 100 replies per second. This matches the rate of connections being sent out which indicates that the NAT is not dropping any connections even with a large number of connections being made. 

The NAT introduced an additional 14 milliseconds in the average reply time over all of the connections. This is likely a result of needing to parse the packets and determine whether or not they should pass through the NAT. 

### Errors 
``` sh 
# No NAT 
Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0

# NAT
Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0
```

While the NAT introduced some latency, it was overall effective in completing connections without any errors or timeouts. This would indicate that the NAT is effectively functional in situations that require create a large amount of connections, with some connections being sent out simultaneously. 

## Without NAT
``` sh
marie@cs352:~$ httperf --server 10.0.0.123 --port 8080 --verbose --rate 100 --num-conn 8000 --timeout 20
httperf --verbose --timeout=20 --client=0/1 --server=10.0.0.123 --port=8080 --uri=/ --rate=100 --send-buffer=4096 --recv-buffer=16384 --num-conns=8000 --num-calls=1
httperf: maximum number of open descriptors = 1048576
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 91.4
reply-rate = 108.8
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
reply-rate = 100.0
Maximum connect burst length: 33

Total: connections 8000 requests 8000 replies 8000 test-duration 79.994 s

Connection rate: 100.0 conn/s (10.0 ms/conn, <=51 concurrent connections)
Connection time [ms]: min 2.6 avg 6.7 max 1152.6 median 4.5 stddev 23.4
Connection time [ms]: connect 1.7
Connection length [replies/conn]: 1.000

Request rate: 100.0 req/s (10.0 ms/req)
Request size [B]: 63.0

Reply rate [replies/s]: min 91.4 avg 100.0 max 108.8 stddev 3.3 (15 samples)
Reply time [ms]: response 4.6 transfer 0.3
Reply size [B]: header 155.0 content 4393.0 footer 0.0 (total 4548.0)
Reply status: 1xx=0 2xx=8000 3xx=0 4xx=0 5xx=0

CPU time [s]: user 42.25 system 32.04 (user 52.8% system 40.1% total 92.9%)
Net I/O: 450.3 KB/s (3.7*10^6 bps)

Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0
```

## With NAT 
``` sh 
marie@cs352:~$ httperf --server 10.0.0.123 --port 8080 --verbose --rate 100 --num-conn 8000 --timeout 20
httperf --verbose --timeout=20 --client=0/1 --server=10.0.0.123 --port=8080 --uri=/ --rate=100 --send-buffer=4096 --recv-buffer=16384 --num-conns=8000 --num-calls=1
httperf: maximum number of open descriptors = 1048576
reply-rate = 99.0
reply-rate = 99.8
reply-rate = 100.4
reply-rate = 100.2
reply-rate = 99.3
reply-rate = 100.8
reply-rate = 99.8
reply-rate = 99.8
reply-rate = 100.2
reply-rate = 100.0
reply-rate = 100.4
reply-rate = 99.8
reply-rate = 100.0
reply-rate = 99.8
reply-rate = 100.0
reply-rate = 100.2
Maximum connect burst length: 1

Total: connections 8000 requests 8000 replies 8000 test-duration 80.027 s

Connection rate: 100.0 conn/s (10.0 ms/conn, <=20 concurrent connections)
Connection time [ms]: min 34.8 avg 48.8 max 1197.5 median 46.5 stddev 23.8
Connection time [ms]: connect 14.5
Connection length [replies/conn]: 1.000

Request rate: 100.0 req/s (10.0 ms/req)
Request size [B]: 63.0

Reply rate [replies/s]: min 99.0 avg 100.0 max 100.8 stddev 0.4 (16 samples)
Reply time [ms]: response 18.3 transfer 15.9
Reply size [B]: header 155.0 content 4393.0 footer 0.0 (total 4548.0)
Reply status: 1xx=0 2xx=8000 3xx=0 4xx=0 5xx=0

CPU time [s]: user 44.96 system 14.15 (user 56.2% system 17.7% total 73.9%)
Net I/O: 450.1 KB/s (3.7*10^6 bps)

Errors: total 0 client-timo 0 socket-timo 0 connrefused 0 connreset 0
Errors: fd-unavail 0 addrunavail 0 ftab-full 0 other 0
```