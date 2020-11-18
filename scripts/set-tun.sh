sudo ip link set dev tun2 up
sudo ip addr add 10.0.0.1/24 dev tun2
sudo ip route add 1.2.3.4 dev tun2
sudo /home/marie/go/bin/control 1 10.0.0.1 0 10.0.2.15 57433
sudo /home/marie/go/bin/control 3 10.0.2.15 0 10.0.0.1 123
sudo /home/marie/go/bin/control 2