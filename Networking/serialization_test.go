package networking

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestSerializeInt(test *testing.T) {
	testint := 42

	buf1 := bytes.NewBuffer(nil)
	SerializeInt(uint64(testint), buf1)

	numasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numasbytes, uint64(testint))
	buf2 := bytes.NewBuffer(numasbytes)

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		test.Error("Int serialization failed")
	}
}

func TestSerializeString(test *testing.T) {
	teststring := "This is a test string"

	buf1 := bytes.NewBuffer(nil)
	var str Netstring
	str.Actstr = teststring
	str.Length = uint64(len(str.Actstr))
	str.Error = false
	SerializeString(str, buf1)

	buf2 := bytes.NewBuffer(nil)
	numasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numasbytes, uint64(len(teststring)))
	buf2.Write(numasbytes)
	buf2.Write([]byte(teststring))
	buf2.Write(make([]byte, 1))

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		test.Error("String serialization failed")
	}
}

func TestSerializeFile(test *testing.T) {
	testfilename := "test.txt"

	var file Netfile
	var str Netstring
	str.Actstr = testfilename
	str.Length = uint64(len(str.Actstr))
	file.Name = str
	testfile, err := os.Open(testfilename)
	if err != nil {
		test.Error("File serialization failed at reading File")
	}
	testfilestat, err := testfile.Stat()
	if err != nil {
		test.Error("File serialization failed at opening Filestat")
	}
	file.Length = uint64(testfilestat.Size())
	file.Actfile = bytes.NewBuffer(nil)
	io.Copy(file.Actfile, testfile)
	buf1 := bytes.NewBuffer(nil)
	testfile.Close()
	file.Error = false
	SerializeFile(file, buf1)

	buf2 := bytes.NewBuffer(nil)
	numasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numasbytes, uint64(len(testfilename)))
	buf2.Write(numasbytes)
	buf2.Write([]byte(testfilename))
	buf2.Write(make([]byte, 1))
	binary.BigEndian.PutUint64(numasbytes, uint64(testfilestat.Size()))
	buf2.Write(numasbytes)
	testfile, err = os.Open(testfilename)
	if err != nil {
		test.Error("File serialization failed at reader File")
	}
	defer testfile.Close()
	io.Copy(buf2, testfile)
	buf2.Write(make([]byte, 1))

	if !bytes.Equal(buf1.Bytes(), buf2.Bytes()) {
		errorstr := fmt.Sprintf("File serialization failed: \n%s %d\n%s %d", buf1.Bytes(), buf1.Len(), buf2.Bytes(), buf2.Len())
		test.Error(errorstr)
	}
}

func TestDeserializeInt(test *testing.T) {
	testint := 42

	numasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numasbytes, uint64(testint))
	buf := bytes.NewBuffer(numasbytes)
	result1 := DeserializeInt(buf)

	buf = bytes.NewBuffer(numasbytes)
	result2 := binary.BigEndian.Uint64(numasbytes)

	if result1 != result2 {
		test.Error("Int deserialization failed")
	}
}

func TestDeserializeString(test *testing.T) {
	teststring := "This is a teststring"

	numasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(numasbytes, uint64(len(teststring)))
	stringasbytes := []byte(teststring)
	buf := bytes.NewBuffer(numasbytes)
	buf.Write(stringasbytes)
	result1 := DeserializeString(buf)

	var result2 Netstring
	result2.Length = binary.BigEndian.Uint64(numasbytes)
	result2.Actstr = string(stringasbytes)

	//if (result1.Length != result2.length) && (result1.Actstr != result2.actstr) {
	if result1 != result2 {
		test.Error("String deserialization failed")
	}
}

func TestDeserializationFile(test *testing.T) {
	testfilename := "test.txt"

	strlenasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(strlenasbytes, uint64(len(testfilename)))
	stringasbytes := []byte(testfilename)
	testfile, err := os.Open(testfilename)
	if err != nil {
		panic(err)
	}
	testfilestat, err := testfile.Stat()
	if err != nil {
		panic(err)
	}
	filelenasbytes := make([]byte, 8)
	binary.BigEndian.PutUint64(filelenasbytes, uint64(testfilestat.Size()))
	buf := bytes.NewBuffer(strlenasbytes)
	buf.Write(stringasbytes)
	buf.Write(make([]byte, 1))
	buf.Write(filelenasbytes)
	io.Copy(buf, testfile)
	testfile.Close()
	buf.Write(make([]byte, 1))
	result1 := DeserializeFile(buf)

	var result2 Netfile
	result2.Name.Length = binary.BigEndian.Uint64(strlenasbytes)
	result2.Name.Actstr = testfilename
	result2.Name.Error = false
	result2.Length = binary.BigEndian.Uint64(filelenasbytes)
	testfile, err = os.Open(testfilename)
	if err != nil {
		panic(err)
	}
	defer testfile.Close()
	result2.Actfile = bytes.NewBuffer(nil)
	io.Copy(result2.Actfile, testfile)
	result2.Error = false

	if !bytes.Equal(result1.Actfile.Bytes(), result2.Actfile.Bytes()) || (result1.Name != result2.Name) || (result1.Length != result2.Length) {
		test.Error("File deserialization failed")
	}
}
