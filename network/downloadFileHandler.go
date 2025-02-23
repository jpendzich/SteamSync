package network

import (
	"os"
	"path/filepath"
)

type DownloadFileHandler struct {
}

func (handler *DownloadFileHandler) Process(conn *Connection) error {
	request := NewDownloadFileRequest()
	response := NewDownloadFileResponse()
	err := conn.ReadPacket(request)
	if err != nil {
		return err
	}
	file, err := os.Open(filepath.Join(request.Game, request.Save))
	if err != nil {
		return err
	}
	defer file.Close()
	response.ErrorCode = 0
	err = conn.WritePacket(response)
	if err != nil {
		return err
	}
	info, err := file.Stat()
	if err != nil {
		return err
	}
	err = conn.WriteBinary(uint64(info.Size()), file)
	if err != nil {
		return err
	}
	return nil
}
