package internal

import (
	"net"
	"os"

	bin "github.com/kelindar/binary"
)

type Connection struct {
	conn    net.Conn
	encoder bin.Encoder
	decoder bin.Decoder
}

func NewConnection(tcpConn *net.Conn) *Connection {
	conn := Connection{}
	conn.conn = *tcpConn
	conn.encoder = *bin.NewEncoder(conn.conn)
	conn.decoder = *bin.NewDecoder(conn.conn)
	return &conn
}

func (conn *Connection) WritePacket(packet interface{}) error {
	err := conn.encoder.Encode(packet)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) WriteFiles(files []os.File) error {
	return nil
}

func (conn *Connection) ReadPacket(packet interface{}) error {
	err := conn.decoder.Decode(packet)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) ReadFiles() ([]os.File, error) {

}
