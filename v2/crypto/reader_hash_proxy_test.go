package ricrypto

import (
	"bytes"
	"fmt"
	"testing"

	"crypto/sha1"
	"hash/fnv"
	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

func TestReaderHashProxy(t *testing.T) {
	b := bytes.NewBufferString("abc")
	h := sha1.New()
	rhp := NewReaderHashProxy(b, h)

	data, err := ioutil.ReadAll(rhp)
	log.PanicIf(err)

	if bytes.Compare(data, []byte{'a', 'b', 'c'}) != 0 {
		t.Fatalf("Data was not read correctly: %v\n", data)
	}

	digestPhrase := fmt.Sprintf("%020x", rhp.Sum())
	if digestPhrase != "a9993e364706816aba3e25717850c26c9cd0d89d" {
		t.Fatalf("hash sum not correct: [%s]", digestPhrase)
	}
}

func TestReaderHash32Proxy(t *testing.T) {
	b := bytes.NewBufferString("abc")
	h := fnv.New32a()
	rhp := NewReaderHash32Proxy(b, h)

	data, err := ioutil.ReadAll(rhp)
	log.PanicIf(err)

	if bytes.Compare(data, []byte{'a', 'b', 'c'}) != 0 {
		t.Fatalf("Data was not read correctly: %v\n", data)
	}

	checksum := rhp.Sum32()
	if checksum != 440920331 {
		t.Fatalf("checksum not correct: (%d)", checksum)
	}
}
