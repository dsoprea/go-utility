package rifs

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

func TestNewBoundedReadWriteSeekerCloser(t *testing.T) {
	sb := NewSeekableBuffer()
	rwsc := ReadWriteSeekNoopCloser(sb)

	brwsc, err := NewBoundedReadWriteSeekCloser(rwsc, 10, 0)
	log.PanicIf(err)

	err = brwsc.Close()
	log.PanicIf(err)
}

func TestNewBoundedReadWriteSeekerCloser_Read(t *testing.T) {
	data := []byte("this is a test; this is a test2")
	sb := NewSeekableBufferWithBytes(data)
	rwsc := ReadWriteSeekNoopCloser(sb)

	brwsc, err := NewBoundedReadWriteSeekCloser(rwsc, 10, 0)
	log.PanicIf(err)

	defer func() {
		err = brwsc.Close()
		log.PanicIf(err)
	}()

	recovered, err := ioutil.ReadAll(brwsc)
	log.PanicIf(err)

	if bytes.Equal(recovered, data[10:]) != true {
		t.Fatalf("Read bytes not correct.")
	}
}

func TestNewBoundedReadWriteSeekerCloser_Write(t *testing.T) {
	sb := NewSeekableBuffer()
	rwsc := ReadWriteSeekNoopCloser(sb)

	brwsc, err := NewBoundedReadWriteSeekCloser(rwsc, 10, 0)
	log.PanicIf(err)

	defer func() {
		err = brwsc.Close()
		log.PanicIf(err)
	}()

	data := []byte("this is a test; this is a test2")

	n, err := brwsc.Write(data)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Written count not correct: (%d)", n)
	}

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	recovered, err := ioutil.ReadAll(sb)
	log.PanicIf(err)

	expectedRecovered := make([]byte, 10+len(data))
	copy(expectedRecovered[10:], data)

	if bytes.Equal(recovered, expectedRecovered) != true {
		fmt.Printf("ACTUAL (%d):\n%v\n\nEXPECTED (%d):\n%v\n", len(recovered), recovered, len(expectedRecovered), expectedRecovered)

		t.Fatalf("Read bytes not correct.")
	}
}

func TestNewBoundedReadWriteSeekerCloser_Seek(t *testing.T) {
	sb := NewSeekableBuffer()
	rwsc := ReadWriteSeekNoopCloser(sb)

	brwsc, err := NewBoundedReadWriteSeekCloser(rwsc, 10, 0)
	log.PanicIf(err)

	defer func() {
		err = brwsc.Close()
		log.PanicIf(err)
	}()

	_, err = brwsc.Seek(10, os.SEEK_SET)
	log.PanicIf(err)

	data := []byte("this is a test; this is a test2")

	n, err := brwsc.Write(data)
	log.PanicIf(err)

	if n != len(data) {
		t.Fatalf("Written count not correct: (%d)", n)
	}

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	recovered, err := ioutil.ReadAll(sb)
	log.PanicIf(err)

	expectedRecovered := make([]byte, 20+len(data))
	copy(expectedRecovered[20:], data)

	if bytes.Equal(recovered, expectedRecovered) != true {
		fmt.Printf("ACTUAL (%d):\n%v\n\nEXPECTED (%d):\n%v\n", len(recovered), recovered, len(expectedRecovered), expectedRecovered)

		t.Fatalf("Read bytes not correct.")
	}
}

func TestNewBoundedReadWriteSeekerCloser_Close(t *testing.T) {
	sb := NewSeekableBuffer()
	rwsc := ReadWriteSeekNoopCloser(sb)

	brwsc, err := NewBoundedReadWriteSeekCloser(rwsc, 10, 0)
	log.PanicIf(err)

	err = brwsc.Close()
	log.PanicIf(err)
}
