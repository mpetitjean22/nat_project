# Configure the tun2 interface
sudo ip link set dev tun2 up
sudo ip addr add 10.0.0.1/24 dev tun2

# the kernal attempts to send reset packets because it does
# not understand what is going out. This rule drops those packets
# so that it does not interfere with the NAT. 
sudo iptables -A OUTPUT -p tcp --tcp-flags RST RST -j DROP
