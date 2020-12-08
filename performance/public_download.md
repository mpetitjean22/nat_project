# wget Performance Results 
## Process
In order to test how the NAT is able to process many packets coming in and out through the same connection, I set up a download test. In this test, we download a large file from the internet. This test is more focused on how the NAT is able to process packets. 

``` sh
# download kernel 
time wget -O /dev/null https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/snapshot/linux-5.10-rc4.tar.gz
```

In order to run the download through the NAT, we must add a route to `https://git.kernel.org`. Since the IP changes, we use the following command: 

``` sh 
# get the IP
KERNEL_IP=$(dig +short ord.git.kernel.org | tail -n 1)

# create route 
sudo ip route add $KERNEL_IP dev tun2
```

This command is also in `scripts/demo.sh` which can be run to acheive the same results. 

## Results 

``` sh 
# No NAT
2020-12-05 13:42:03 (3.95 MB/s) - ‘/dev/null’ saved [185282209]

# NAT
2020-12-05 13:39:30 (3.51 MB/s) - ‘/dev/null’ saved [185282209]

```

We can see that after downloading 176.70 MB from `https://git.kernel.org`, going through the NAT has an average rate of download of 3.51 MB/s whereas without the NAT has an avergae rate of 3.95 MB/s. We also know that there is a difference of 5 seconds in the download time, which means that over 176.70 MB, the NAT introduces 28.2 milliseconds per MB, which is equivalent to only 0.282 ns per byte, which is a very small amount of latency. 


# Data
## No NAT
``` sh
marie@cs352:~/nat_project$ time wget -O /dev/null https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/snapshot/linux-5.10-rc4.tar.gz
--2020-12-05 13:41:18--  https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/snapshot/linux-5.10-rc4.tar.gz
Resolving git.kernel.org (git.kernel.org)... 147.75.58.133, 2604:1380:4020:600::1
Connecting to git.kernel.org (git.kernel.org)|147.75.58.133|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: unspecified [application/x-gzip]
Saving to: ‘/dev/null’

/dev/null                             [           <=>                                               ] 176.70M  5.24MB/s    in 45s

2020-12-05 13:42:03 (3.95 MB/s) - ‘/dev/null’ saved [185282209]


real	0m44.917s
user	0m0.312s
sys	0m2.707s
```

## NAT 
``` sh
marie@cs352:~/nat_project$ time wget -O /dev/null https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/snapshot/linux-5.10-rc4.tar.gz
--2020-12-05 13:38:39--  https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/snapshot/linux-5.10-rc4.tar.gz
Resolving git.kernel.org (git.kernel.org)... 147.75.58.133, 2604:1380:4020:600::1
Connecting to git.kernel.org (git.kernel.org)|147.75.58.133|:443... connected.
HTTP request sent, awaiting response... 200 OK
Length: unspecified [application/x-gzip]
Saving to: ‘/dev/null’

/dev/null                             [          <=>                                                ] 176.70M  3.87MB/s    in 50s

2020-12-05 13:39:30 (3.51 MB/s) - ‘/dev/null’ saved [185282209]


real	0m50.615s
user	0m1.236s
sys	0m0.000s
```