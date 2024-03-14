package utils

type PFile struct {
	name    PString
	length  PInt
	payload []byte
}

type PString struct {
	length  PInt
	payload string
}

type PInt struct {
	payload uint64
}
