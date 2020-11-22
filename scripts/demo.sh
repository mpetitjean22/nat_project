# Add a route to route all packets from google and wikipedia to 
# the tun2 interface. This allows us to use lynx in a demo. Lynx
# does not support binding to a particular interface so these 
# rules are necessary. 

sudo ip route add 172.217.0.0/16 dev tun2
sudo ip route add 142.250.0.0/16 dev tun2
sudo ip route add 208.80.0.0/16 dev tun2
