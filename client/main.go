package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/songgao/water"
)

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.3"

func main() {
	ip := *flag.String("ip", "", "IP address of the proxy server")
	host := *flag.String("host", "", "Domain of the server to imitate")

	flag.Parse()

	if ip == "" || host == "" {
		flag.Usage()
		os.Exit(1)
	}

	gateway, err := NewGateway(ip, host)
	if err != nil {
		log.Fatal(err.Error())
	}

	adapter, err := NewAdapter()
	if err != nil {
		log.Fatal(err.Error())
	}

	go Relay(adapter, gateway)
	go Relay(gateway, adapter)
}

type Transport interface {
	TX(buffer []byte) error
	RX() ([]byte, error)
}

func Relay(src, dst Transport) {
	for {
		buffer, err := src.RX()
		if err != nil {
			log.Println("Unable Read: ", err.Error())
			break
		}

		err = dst.TX(buffer)
		if err != nil {
			log.Println("Unable Write: ", err.Error())
			break
		}
	}
}

type Gateway struct {
	Conn *websocket.Conn
}

func NewGateway(ip, host string) (*Gateway, error) {
	headers := http.Header{}
	headers.Set("Host", host)
	headers.Set("User-Agent", USER_AGENT)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(fmt.Sprintf("http://%s/", ip), headers)
	if err != nil {
		return nil, err
	}

	return &Gateway{
		Conn: conn,
	}, nil
}

func (g *Gateway) TX(buffer []byte) error {
	return g.Conn.WriteMessage(websocket.BinaryMessage, buffer)
}

func (g *Gateway) RX() ([]byte, error) {
	type_, buffer, err := g.Conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}

	if type_ != websocket.BinaryMessage {
		return []byte{}, errors.New("Expected a binary message")
	}

	return buffer, nil
}

type Adapter struct {
	Interface *water.Interface
}

func NewAdapter() (*Adapter, error) {
	config := water.Config{
		DeviceType: water.TUN,
	}
	config.Name = "wrapper"

	ifce, err := water.New(config)
	if err != nil {
		return nil, err
	}

	return &Adapter{
		Interface: ifce,
	}, nil
}

func (a *Adapter) TX(buffer []byte) error {
	_, err := a.Interface.Write(buffer)
	return err
}

func (a *Adapter) RX() ([]byte, error) {
	packet := make([]byte, 2000)
	n, err := a.Interface.Read(packet)
	return packet[:n], err
}
