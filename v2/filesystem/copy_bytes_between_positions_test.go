package rifs

import (
	"bytes"
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestCopyBytesBetweenPositions(t *testing.T) {
	sb := NewSeekableBuffer()

	testData1 := []byte("abcde")

	_, err := sb.Write(testData1)
	log.PanicIf(err)

	_, err = sb.Seek(20, os.SEEK_SET)
	log.PanicIf(err)

	testData2 := []byte("fghij")

	_, err = sb.Write(testData2)
	log.PanicIf(err)

	expected1 := make([]byte, 25)
	copy(expected1[0:], testData1)
	copy(expected1[20:], testData2)

	if bytes.Compare(sb.Bytes(), expected1) != 0 {
		t.Fatalf("Initial insert not correctly laid-out: %v", sb.Bytes())
	}

	expectedByteCount1 := len(testData2)
	n, err := CopyBytesBetweenPositions(sb, 20, 10, expectedByteCount1)
	log.PanicIf(err)

	if n != expectedByteCount1 {
		t.Fatalf("copied bytes not correct: (%d) != (%d)", n, expectedByteCount1)
	}

	expected2 := make([]byte, 25)
	copy(expected2[0:], testData1)
	copy(expected2[20:], testData2)
	copy(expected2[10:], testData2)

	if bytes.Compare(sb.Bytes(), expected2) != 0 {
		t.Fatalf("Copy result not correctly laid-out: %v", sb.Bytes())
	}
}
