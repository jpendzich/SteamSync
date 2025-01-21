package network

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

type UploadFileHandler struct {
}

func (handler *UploadFileHandler) Process(conn *Connection) error {
	request := NewUploadFileRequest()
	response := NewUploadFileResponse()
	err := conn.ReadPacket(&request)
	if err != nil {
		return err
	}
	_, err = os.Stat(request.Game)
	if errors.Is(err, fs.ErrNotExist) {
		os.Mkdir(request.Game, os.ModePerm)
	} else if err != nil {
		return err
	}
	file, err := os.Create(path.Join(request.Game, request.Save))
	if err != nil {
		return err
	}
	defer file.Close()
	err = conn.ReadBinary(file)
	if err != nil {
		return err
	}

	response.ErrorCode = ErrorOk
	err = conn.WritePacket(response)
	if err != nil {
		return err
	}
	return nil
}
