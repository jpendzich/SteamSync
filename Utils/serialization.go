package utils

import (
	"encoding/binary"
)

func SerializeFile(packet PFile) []byte {
	bytes := SerializeString(packet.name)
	bytes = append(bytes, SerializeInt(packet.length)...)
	bytes = append(bytes, packet.payload...)
	return bytes
}

func SerializeString(packet PString) []byte {
	bytes := SerializeInt(packet.length)
	bytes = append(bytes, []byte(packet.payload)...)
	return bytes
}

func SerializeInt(packet PInt) []byte {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, packet.payload)
	return bytes
}

func DeserializeFile(bytes *[]byte) PFile {
	var packet PFile
	newbytes := *bytes
	packet.name = DeserializeString(&newbytes)
	packet.length = DeserializeInt(&newbytes)
	packet.payload = newbytes[:packet.length.payload]
	newbytes = newbytes[packet.length.payload:]
	*bytes = newbytes
	return packet
}

func DeserializeString(bytes *[]byte) PString {
	var packet PString
	newbytes := *bytes
	packet.length = DeserializeInt(&newbytes)
	packet.payload = string(newbytes[:packet.length.payload])
	newbytes = newbytes[packet.length.payload:]
	*bytes = newbytes
	return packet
}

func DeserializeInt(bytes *[]byte) PInt {
	var packet PInt
	newbytes := *bytes
	packet.payload = binary.BigEndian.Uint64(newbytes[:8])
	newbytes = newbytes[8:]
	*bytes = newbytes
	return packet
}
