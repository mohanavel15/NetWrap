package pkg

import (
	"errors"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrReadWSConn   = errors.New("ErrReadWSConn")
	ErrWrongMsgType = errors.New("ErrWrongMsgType")
)

type WrappedWsConn struct {
	conn *websocket.Conn
}

func NewWrappedWsConn(conn *websocket.Conn) WrappedWsConn {
	return WrappedWsConn{conn: conn}
}

func (tc WrappedWsConn) Read(b []byte) (int, error) {
	t, src, err := tc.conn.ReadMessage()
	if err != nil {
		return 0, ErrReadWSConn
	}

	if t != websocket.BinaryMessage {
		return 0, ErrWrongMsgType
	}

	n := copy(b, src)
	return n, nil
}

func (tc WrappedWsConn) Write(b []byte) (int, error) {
	err := tc.conn.WriteMessage(websocket.BinaryMessage, b)
	return len(b), err
}

func (tc WrappedWsConn) Close() error {
	return tc.conn.Close()
}

func (tc WrappedWsConn) LocalAddr() net.Addr {
	return nil
}

func (tc WrappedWsConn) RemoteAddr() net.Addr {
	return nil
}

func (tc WrappedWsConn) SetDeadline(t time.Time) error {
	err1 := tc.conn.SetReadDeadline(t)
	err2 := tc.conn.SetWriteDeadline(t)
	return errors.Join(err1, err2)
}

func (tc WrappedWsConn) SetReadDeadline(t time.Time) error {
	return tc.conn.SetReadDeadline(t)
}

func (tc WrappedWsConn) SetWriteDeadline(t time.Time) error {
	return tc.conn.SetWriteDeadline(t)
}
