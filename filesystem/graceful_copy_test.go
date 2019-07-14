package rifs

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestGracefulCopy__Defaults(t *testing.T) {
	data := []byte("test bytes")
	sbFrom := NewSeekableBufferWithBytes(data)
	sbTo := NewSeekableBuffer()

	n, err := GracefulCopy(sbTo, sbFrom, nil)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Count of copied bytes not correct: (%d)", n)
	}

	copiedBytes := sbTo.Bytes()

	if bytes.Equal(copiedBytes, data) != true {
		t.Fatalf("Copied bytes not correct.")
	}
}

func TestGracefulCopy__ShortBuffer(t *testing.T) {
	data := []byte("test bytes")
	sbFrom := NewSeekableBufferWithBytes(data)
	sbTo := NewSeekableBuffer()

	buffer := make([]byte, 3)
	n, err := GracefulCopy(sbTo, sbFrom, buffer)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Count of copied bytes not correct: (%d)", n)
	}

	copiedBytes := sbTo.Bytes()

	if bytes.Equal(copiedBytes, data) != true {
		t.Fatalf("Copied bytes not correct.")
	}
}

func TestGracefulCopy__LongBuffer(t *testing.T) {
	data := []byte("test bytes")
	sbFrom := NewSeekableBufferWithBytes(data)
	sbTo := NewSeekableBuffer()

	buffer := make([]byte, 50)
	n, err := GracefulCopy(sbTo, sbFrom, buffer)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Count of copied bytes not correct: (%d)", n)
	}

	copiedBytes := sbTo.Bytes()

	if bytes.Equal(copiedBytes, data) != true {
		t.Fatalf("Copied bytes not correct.")
	}
}

func TestGracefulCopy__BigData(t *testing.T) {
	data := []byte(strings.Repeat("test bytes", 1024*1024))
	sbFrom := NewSeekableBufferWithBytes(data)
	sbTo := NewSeekableBuffer()

	n, err := GracefulCopy(sbTo, sbFrom, nil)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Count of copied bytes not correct: (%d)", n)
	}

	copiedBytes := sbTo.Bytes()

	if bytes.Equal(copiedBytes, data) != true {
		t.Fatalf("Copied bytes not correct.")
	}
}

// shortWriter will not write the complete slice for data larger than a certain
// threshold.
type shortWriter struct {
	w io.Writer
}

func (sh shortWriter) Write(buffer []byte) (n int, err error) {
	if len(buffer) > 50 {
		buffer = buffer[:50]
	}

	n, err = sh.w.Write(buffer)
	if err != nil {
		panic(err)
	}

	return n, nil
}

func TestGracefulCopy__ShortWrite(t *testing.T) {
	data := []byte(strings.Repeat("test bytes", 1024*1024))
	sbFrom := NewSeekableBufferWithBytes(data)

	sbTo := NewSeekableBuffer()
	sw := &shortWriter{
		w: sbTo,
	}

	n, err := GracefulCopy(sw, sbFrom, nil)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Count of copied bytes not correct: (%d)", n)
	}

	copiedBytes := sbTo.Bytes()

	if bytes.Equal(copiedBytes, data) != true {
		t.Fatalf("Copied bytes not correct.")
	}
}

// uncleanEofReader will return EOF with nonzero read bytes.
type uncleanEofReader struct {
	r io.Reader
}

func (uer uncleanEofReader) Read(buffer []byte) (n int, err error) {
	n, err = uer.r.Read(buffer)
	if err == io.EOF {
		return 0, err
	} else if err != nil {
		panic(err)
	}

	// Return EOF every time. The only time it should make a difference is when
	// (n == 0).
	return n, io.EOF
}

func TestGracefulCopy__UncleanEof(t *testing.T) {
	data := []byte(strings.Repeat("test bytes", 1024*1024))
	sbFrom := NewSeekableBufferWithBytes(data)
	uer := &uncleanEofReader{
		r: sbFrom,
	}

	sbTo := NewSeekableBuffer()
	n, err := GracefulCopy(sbTo, uer, nil)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Count of copied bytes not correct: (%d)", n)
	}

	copiedBytes := sbTo.Bytes()

	if bytes.Equal(copiedBytes, data) != true {
		t.Fatalf("Copied bytes not correct.")
	}
}
