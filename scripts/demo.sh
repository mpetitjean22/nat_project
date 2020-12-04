# Add a route to route all packets from google and wikipedia to 
# the tun2 interface. This allows us to use lynx in a demo. Lynx
# does not support binding to a particular interface so these 
# rules are necessary. 

sudo ip route add 172.217.0.0/16 dev tun2
sudo ip route add 208.80.0.0/16 dev tun2

KERNEL_IP=$(dig +short ord.git.kernel.org | tail -n 1)

sudo ip route add $KERNEL_IP dev tun2

echo Make sure $KERNEL_IP is in mappings after download. DNS could have changed between routing and running lynx.
