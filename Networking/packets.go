package networking

import (
	"bytes"
	"io"
)

type Packet interface {
	Serialize(dest io.Writer) error
	Deserialize(src io.Reader) error
}

type Request struct {
	Version   uint
	ErrorCode uint
}

type Response struct {
	Version   uint
	ErrorCode uint
}

type UploadFileRequest struct {
	Base Request
	Name string
	Data bytes.Buffer
}

type UploadFileResponse struct {
	Base Response
}

type DownloadFileRequest struct {
	Base Request
	Game string
	Save string
}

type DownloadFileResponse struct {
	Base Response
	Data bytes.Buffer
}
