package network

import (
	"io"
	"log"
)

type RequestSender struct {
	conn *Connection
}

func NewRequestSender(conn *Connection) *RequestSender {
	return &RequestSender{
		conn: conn,
	}
}

func (sender *RequestSender) SendRequest(packetType uint, request interface{}, response interface{}) error {
	sender.conn.WritePacketType(packetType)
	err := sender.conn.WritePacket(request)
	if err != nil {
		return err
	}
	err = sender.conn.ReadPacket(response)
	if err != nil {
		return err
	}
	return nil
}

func (sender *RequestSender) SendRequestWriteBinary(packetType uint, request interface{}, response interface{}, binary io.Reader, size uint64) error {
	sender.conn.WritePacketType(packetType)
	err := sender.conn.WritePacket(request)
	if err != nil {
		return err
	}
	err = sender.conn.WriteBinary(size, binary)
	if err != nil {
		return err
	}
	err = sender.conn.ReadPacket(response)
	if err != nil {
		return err
	}
	return nil
}

func (sender *RequestSender) SendRequestReadBinary(packetType uint, request interface{}, response interface{}, binary io.Writer) error {
	sender.conn.WritePacketType(packetType)
	log.Println(request)
	err := sender.conn.WritePacket(request)
	if err != nil {
		return err
	}
	err = sender.conn.ReadPacket(response)
	if err != nil {
		return err
	}
	err = sender.conn.ReadBinary(binary)
	if err != nil {
		return err
	}
	return nil
}
