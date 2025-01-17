package networking

import (
	"encoding/binary"
	"io"
)

func SerializeUint(value uint, dest io.Writer) error {
	buf := make([]byte, 8) // Uint is always 64 bit in this protocol
	binary.BigEndian.PutUint64(buf, uint64(value))
	_, err := dest.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func DeserializeUint(src io.Reader) (uint, error) {
	buf := make([]byte, 8)
	_, err := src.Read(buf)
	if err != nil {
		return 0, err
	}
	return uint(binary.BigEndian.Uint64(buf)), nil
}

func SerializeString(str string, dest io.Writer) error {
	length := make([]byte, 8)
	binary.BigEndian.PutUint64(length, uint64(len(str)))
	_, err := dest.Write(length)
	if err != nil {
		return err
	}
	_, err = dest.Write([]byte(str))
	if err != nil {
		return err
	}
	return nil
}

func DeserializeString(src io.Reader) (string, error) {
	length, err := DeserializeUint(src)
	if err != nil {
		return "", err
	}
	nameBuf := make([]byte, length)
	_, err = src.Read(nameBuf)
	if err != nil {
		return "", err
	}
	return string(nameBuf), nil
}

func (request *Request) Serialize(dest io.Writer) error {
	err := SerializeUint(request.Version, dest)
	if err != nil {
		return err
	}
	err = SerializeUint(request.ErrorCode, dest)
	if err != nil {
		return err
	}
	return nil
}

func (request *Request) Deserialize(src io.Reader) error {
	var err error
	request.Version, err = DeserializeUint(src)
	if err != nil {
		return err
	}
	request.ErrorCode, err = DeserializeUint(src)
	if err != nil {
		return err
	}
	return nil
}

func (response *Response) Serialize(dest io.Writer) error {
	err := SerializeUint(response.Version, dest)
	if err != nil {
		return err
	}
	err = SerializeUint(response.ErrorCode, dest)
	if err != nil {
		return err
	}
	return nil
}

func (response *Response) Deserialize(src io.Reader) error {
	var err error
	response.Version, err = DeserializeUint(src)
	if err != nil {
		return err
	}
	response.ErrorCode, err = DeserializeUint(src)
	if err != nil {
		return err
	}
	return nil
}

func (request *UploadFileRequest) Serialize(dest io.Writer) error {
	err := request.Base.Serialize(dest)
	if err != nil {
		return err
	}
	length := make([]byte, 8)
	SerializeString(request.Name, dest)
	binary.BigEndian.PutUint64(length, uint64(request.Data.Len()))
	dest.Write(length)
	_, err = io.Copy(dest, &request.Data)
	if err != nil {
		return err
	}
	return nil
}

func (request *UploadFileRequest) Deserialize(src io.Reader) error {
	err := request.Base.Deserialize(src)
	if err != nil {
		return err
	}
	length := make([]byte, 8)
	request.Name, err = DeserializeString(src)
	if err != nil {
		return err
	}
	src.Read(length)
	_, err = io.CopyN(&request.Data, src, int64(binary.BigEndian.Uint64(length)))
	if err != nil {
		return err
	}
	return nil
}

func (response *UploadFileResponse) Serialize(dest io.Writer) error {
	err := response.Base.Serialize(dest)
	if err != nil {
		return err
	}
	return nil
}

func (response *UploadFileResponse) Deserialize(src io.Reader) error {
	err := response.Base.Deserialize(src)
	if err != nil {
		return err
	}
	return nil
}

func (request *DownloadFileRequest) Serialize(dest io.Writer) error {
	err := request.Base.Serialize(dest)
	if err != nil {
		return err
	}
}

func (request *DownloadFileRequest) Deserialize(src io.Reader) error {

}

func (response *DownloadFileResponse) Serialize(dest io.Writer) error {

}

func (response *DownloadFileResponse) Deserialize(src io.Reader) error {

}
