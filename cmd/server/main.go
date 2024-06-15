package main

import (
	"NetWrap/pkg"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func main() {
	port := *flag.String("port", "5000", "Port number for the server to run on")
	flag.Parse()

	adapter, err := pkg.NewAdapter("exit_node")
	if err != nil {
		log.Fatal("Error:", err.Error())
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: &HttpHandler{Adp: adapter},
	}

	server.ListenAndServe()
}

var upgrader = websocket.Upgrader{}

type HttpHandler struct {
	Adp *pkg.Adapter
}

func (h *HttpHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}

	node := pkg.NewNodeFromConn(conn)
	go pkg.Relay(node, h.Adp)
	pkg.Relay(h.Adp, node)
}
