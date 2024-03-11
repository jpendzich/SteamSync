package utils

import "encoding/binary"

func SerializeFile(packet PFile) []byte {

}

func SerializeString(packet PString) []byte {

}

func SerializeInt(packet PInt) []byte {

}

func DeserializeFile(bytes []byte) PFile {
	var packet PFile
	packet.length = DeserializeInt(bytes[8:]).payload
	bytes = bytes[:8]
	packet.name = DeserializeString(bytes).payload
	return packet
}

func DeserializeString(bytes []byte) PString {
	var packet PString
	packet.length = DeserializeInt(bytes[8:]).payload
	bytes = bytes[:8]
	packet.payload = string(bytes[:packet.length])
	bytes = bytes[:packet.length]
	return packet
}

func DeserializeInt(bytes []byte) PInt {
	var packet PInt
	binary.BigEndian.PutUint64(bytes, uint64(packet.payload))
	return packet
}
