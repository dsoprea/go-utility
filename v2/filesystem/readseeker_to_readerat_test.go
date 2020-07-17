package rifs

import (
	"bytes"
	"io"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestNewReadSeekerToReaderAt(t *testing.T) {
	b := []byte{}
	br := bytes.NewReader(b)
	rstra := NewReadSeekerToReaderAt(br)

	if rstra.rs != io.ReadSeeker(br) {
		t.Fatalf("Inner ReadSeeker instance is not correct.")
	}
}

func TestReadSeekerToReaderAt_ReadAt(t *testing.T) {
	b := []byte{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	br := bytes.NewReader(b)
	rstra := NewReadSeekerToReaderAt(br)

	offsets := []int{15, 7, 12, 0, 3, 17}
	lengths := []int{5, 5, 5, 5, 5, 3}

	for i, dataOffset := range offsets {
		dataLength := lengths[i]

		initialOffset, err := br.Seek(0, io.SeekCurrent)
		log.PanicIf(err)

		if initialOffset != 0 {
			t.Fatalf("Before-ReadAt offset not correct: (%d)", initialOffset)
		}

		fragment := make([]byte, dataLength)

		_, err = rstra.ReadAt(fragment, int64(dataOffset))
		log.PanicIf(err)

		if bytes.Equal(fragment, b[dataOffset:dataOffset+dataLength]) != true {
			t.Fatalf("Seek (1) returned incorrect bytes: %v\n", fragment)
		}

		currentOffset, err := br.Seek(0, io.SeekCurrent)
		log.PanicIf(err)

		if currentOffset != initialOffset {
			t.Fatalf("After-ReadAt offset not correct: (%d) != (%d)", currentOffset, initialOffset)
		}
	}
}
