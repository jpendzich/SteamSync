package network

import (
	"io"
	"net"

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

func (conn *Connection) WritePacketType(command uint) {
	conn.encoder.WriteUint32(uint32(command))
}

func (conn *Connection) WritePacket(packet interface{}) error {
	err := conn.encoder.Encode(packet)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) WriteBinary(size uint64, binary io.Reader) error {
	conn.encoder.WriteUint64(size)
	_, err := io.Copy(conn.conn, binary)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) ReadPacketType() (uint, error) {
	command, err := conn.decoder.ReadUint32()
	if err != nil {
		return 0, err
	}
	return uint(command), nil
}

func (conn *Connection) ReadPacket(packet interface{}) error {
	err := conn.decoder.Decode(packet)
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) ReadBinary(binary io.Writer) error {
	length, err := conn.decoder.ReadUint64()
	if err != nil {
		return err
	}

	_, err = io.CopyN(binary, &conn.decoder, int64(length))
	if err != nil {
		return err
	}
	return nil
}
