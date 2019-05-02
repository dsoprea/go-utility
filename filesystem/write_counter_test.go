package rifs

import (
	"bytes"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestWriteCounter_Write(t *testing.T) {
	s := "testing string"
	b := bytes.NewBufferString(s)

	wc := NewWriteCounter(b)

	n, err := wc.Write([]byte(s))
	log.PanicIf(err)

	writtenBytes := b.Bytes()[:n]

	if n != len(s) {
		t.Fatalf("Count of bytes written not correct: (%d)", n)
	} else if bytes.Compare(writtenBytes, []byte(s)) != 0 {
		t.Fatalf("Written data not correct:\nACTUAL:\n%v\nEXPECTED:\n%v", writtenBytes, []byte(s))
	}

	if wc.Count() != len(s) {
		t.Fatalf("Counter not correct: (%d)", wc.Count())
	}

	wc.Reset()

	if wc.Count() != 0 {
		t.Fatalf("Count after reset not zero: (%d)", wc.Count())
	}
}
