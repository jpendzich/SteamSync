package internal

type Packet interface {
}

type Request struct {
	Version int
	Command uint
}

type Response struct {
	Version   int
	ErrorCode int
	Command   uint
}

type UploadFileRequest struct {
	Request
	Game     string
	Save     string
	NumFiles uint
}

type UploadFileResponse struct {
	Response
}

type DownloadFileRequest struct {
	Request
	Game string
	Save string
}

type DownloadFileResponse struct {
	Response
	NumFiles uint
}
