package main

import (
	"NetWrap/pkg"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

func main() {
	ip := flag.String("ip", "", "IP address of the proxy server")
	host := flag.String("host", "", "Domain of the server to imitate")
	flag.Parse()

	if *ip == "" || *host == "" {
		flag.Usage()
		os.Exit(1)
	}

	adapter, err := pkg.NewAdapter("netwrap")
	if err != nil {
		log.Fatal(err.Error())
	}

	headers := http.Header{}
	headers.Set("Host", *host)
	headers.Set("User-Agent", pkg.USER_AGENT)

	dialer := websocket.DefaultDialer
	wsconn, _, err := dialer.Dial(fmt.Sprintf("ws://%s/", *ip), headers)
	if err != nil {
		fmt.Println("Unable to connect to gateway!:", err.Error())
		os.Exit(1)
	}

	wconn := pkg.NewWrappedWsConn(wsconn)
	config := &tls.Config{
		ServerName:         "example.com",
		InsecureSkipVerify: false,
		Certificates:       []tls.Certificate{},
		RootCAs:            nil,
		MinVersion:         tls.VersionTLS12,
		MaxVersion:         tls.VersionTLS13,
	}

	conn := tls.Client(wconn, config)

	// gateway, err := pkg.NewNode(*ip, *host)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// tls.Client(gateway.Conn, nil)

	go io.Copy(adapter.Interface, conn)
	io.Copy(conn, adapter.Interface)

	// go pkg.Relay(adapter, gateway)
	// pkg.Relay(gateway, adapter)
}
