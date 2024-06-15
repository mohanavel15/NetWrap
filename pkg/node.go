package pkg

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Node struct {
	Conn *websocket.Conn
}

func NewNode(ip, host string) (*Node, error) {
	headers := http.Header{}
	headers.Set("Host", host)
	headers.Set("User-Agent", USER_AGENT)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(fmt.Sprintf("ws://%s/", ip), headers)
	if err != nil {
		return nil, err
	}

	return &Node{
		Conn: conn,
	}, nil
}

func NewNodeFromConn(conn *websocket.Conn) *Node {
	return &Node{
		Conn: conn,
	}
}

func (g *Node) TX(buffer []byte) error {
	return g.Conn.WriteMessage(websocket.BinaryMessage, buffer)
}

func (g *Node) RX() ([]byte, error) {
	type_, buffer, err := g.Conn.ReadMessage()
	if err != nil {
		return []byte{}, err
	}

	if type_ != websocket.BinaryMessage {
		return []byte{}, errors.New("Expected a binary message")
	}

	return buffer, nil
}
