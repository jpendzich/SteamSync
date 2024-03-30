package utils

type PFile struct {
	lengthpacket  PInt
	name          PString
	lengthpayload PInt
	payload       []byte
}

type PString struct {
	length  PInt
	payload string
}

type PInt struct {
	payload uint64
}
