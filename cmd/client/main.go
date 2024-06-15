package main

import (
	"NetWrap/pkg"
	"flag"
	"log"
	"os"
)

func main() {
	ip := *flag.String("ip", "127.0.0.1:5000", "IP address of the proxy server")
	host := *flag.String("host", "www.netflix.com", "Domain of the server to imitate")

	if ip == "" || host == "" {
		flag.Usage()
		os.Exit(1)
	}

	gateway, err := pkg.NewNode(ip, host)
	if err != nil {
		log.Fatal(err.Error())
	}

	adapter, err := pkg.NewAdapter("start_node")
	if err != nil {
		log.Fatal(err.Error())
	}

	go pkg.Relay(adapter, gateway)
	pkg.Relay(gateway, adapter)
}
