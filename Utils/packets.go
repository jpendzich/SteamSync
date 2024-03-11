package utils

type PFile struct {
	name    string
	length  int64
	payload []byte
}

type PString struct {
	length  int64
	payload string
}

type PInt struct {
	payload int64
}
