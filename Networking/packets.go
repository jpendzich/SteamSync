package networking

import (
	"bytes"
	"io"
)

type Packet interface {
	Serialize(dest io.Writer) error
	Deserialize(src io.Reader) error
}

type UploadFileRequest struct {
	name string
	data bytes.Buffer
}

type UploadFileResponse struct {
	errorCode uint
}

type DownloadFileRequest struct {
	name string
	data bytes.Buffer
}

type DownloadFileResponse struct {
	errorCode uint
}
