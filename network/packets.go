package network

type Packet interface {
}

type Request struct {
	Version    int
	PacketType uint
}

type Response struct {
	Version    int
	ErrorCode  int
	PacketType uint
}

type UploadFileRequest struct {
	Request
	Game string
	Save string
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
}

func NewUploadFileRequest() *UploadFileRequest {
	return &UploadFileRequest{
		Request: Request{Version: 1, PacketType: 1},
	}
}

func NewUploadFileResponse() *UploadFileResponse {
	return &UploadFileResponse{
		Response: Response{Version: 1, PacketType: 2},
	}
}

func NewDownloadFileRequest() *DownloadFileRequest {
	return &DownloadFileRequest{
		Request: Request{Version: 1, PacketType: 3},
	}
}

func NewDownloadFileResponse() *DownloadFileResponse {
	return &DownloadFileResponse{
		Response: Response{Version: 1, PacketType: 4},
	}
}
