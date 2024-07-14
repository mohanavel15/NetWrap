package main

import (
	"NetWrap/pkg"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

func main() {
	port := flag.String("port", "5000", "Port number for the server to run on")
	flag.Parse()

	adapter, err := pkg.NewAdapter("tunnel")
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	// err = adapter.SetUpFd()
	// fmt.Println(err)
	// err = adapter.SetIP(net.IPv4(10, 10, 1, 1))
	// fmt.Println(err)
	// err = adapter.SetUp()
	// fmt.Println(err)
	configureIPTables()

	conns := make(chan net.Conn, 10)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", *port),
		Handler: &HttpHandler{conns: conns},
	}

	go server.ListenAndServe()

	for {

		conn := <-conns
		go io.Copy(adapter.Interface, conn)
		go io.Copy(conn, adapter.Interface)
	}
}

func configureIPTables() error {
	// iptables -t nat -A POSTROUTING -j MASQUERADE
	cmd := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-j", "MASQURADE")
	return cmd.Run()
}

var upgrader = websocket.Upgrader{}

type HttpHandler struct {
	conns chan net.Conn
}

func (h *HttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}

	wconn := pkg.NewWrappedWsConn(conn)
	h.conns <- wconn
	// node := pkg.NewNodeFromConn(conn)
	// go pkg.Relay(node, h.Adp)
	// pkg.Relay(h.Adp, node)
}
