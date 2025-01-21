package network

var registry map[uint]Handler

type Handler interface {
	Process(conn *Connection) error
}

func RegisterHandler() {
	registry = make(map[uint]Handler)
	registry[1] = &UploadFileHandler{}
	registry[3] = &DownloadFileHandler{}
}

func HandlePacket(conn *Connection, packetType uint) error {
	return registry[packetType].Process(conn)
}
