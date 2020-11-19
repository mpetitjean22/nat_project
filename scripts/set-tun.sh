sudo ip link set dev tun2 up
sudo ip addr add 10.0.0.1/24 dev tun2
sudo iptables -A OUTPUT -p tcp --tcp-flags RST RST -j DROP
