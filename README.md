# NetWrap
Simple Network Tunnel

It's not that great yet.

# Server set up(on linux):
```
go run cmd/server/main.go
```
or
```
make server
cd build
./server
```

### Enable ip forwarding:
Add or uncomment line `net.ipv4.ip_forward=1` in `/etc/sysctl.conf`
```
$ sysctl -p
```
### Configure the network interface and the system:
```
$ ip addr add 192.168.1.1/24 dev {interface}
$ ip link set dev {interface} up
$ iptables -t nat -A POSTROUTING -j MASQUERADE
```

# Client set up(on linux):
```
go run cmd/client/main.go
```
or
```
make client
cd build
./client
```

### Configure the network interface and the system.
```
$ ip addr add 192.168.1.{2-254}/24 dev {interface}
$ ip link set dev {interface} up
$ ip route add default via 192.168.1.1 dev {interface}
```