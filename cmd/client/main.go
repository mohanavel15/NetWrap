package main

import (
	"NetWrap/pkg"
	"flag"
	"log"
	"os"
)

func main() {
	ip := flag.String("ip", "", "IP address of the proxy server")
	host := flag.String("host", "", "Domain of the server to imitate")
	flag.Parse()

	if *ip == "" || *host == "" {
		flag.Usage()
		os.Exit(1)
	}

	gateway, err := pkg.NewNode(*ip, *host)
	if err != nil {
		log.Fatal(err.Error())
	}

	adapter, err := pkg.NewAdapter("netwrap")
	if err != nil {
		log.Fatal(err.Error())
	}

	go pkg.Relay(adapter, gateway)
	pkg.Relay(gateway, adapter)
}
