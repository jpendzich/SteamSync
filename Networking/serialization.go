package networking

import (
	"encoding/binary"
	"io"
)

func (request *UploadFileRequest) Serialize(dest io.Writer) error {
	const packetName = "UploadFileRequest"
	length := make([]byte, 8)
	binary.BigEndian.PutUint64(length, uint64(len(packetName)))
	dest.Write(length)
	dest.Write([]byte(packetName))
	binary.BigEndian.PutUint64(length, uint64(len(request.name)))
	dest.Write(length)
	dest.Write([]byte(request.name))
	binary.BigEndian.PutUint64(length, uint64(request.data.Len()))
	dest.Write(length)
	io.Copy(dest, &request.data)
	return nil
}

func (response *UploadFileResponse) Deserialize(src io.Reader) error {
	buf := make([]byte, 4)
	_, err := src.Read(buf)
	if err != nil {
		return err
	}
	response.errorCode = uint(binary.BigEndian.Uint32(buf))
	return nil
}
