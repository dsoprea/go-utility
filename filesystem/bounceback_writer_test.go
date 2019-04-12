package rifs

import (
	"bytes"
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestBouncebackWriter(t *testing.T) {
	sb := NewSeekableBuffer()

	bw, err := NewBouncebackWriter(sb)
	log.PanicIf(err)

	input := []byte("testinput")

	for i := 0; i < len(input); i++ {
		_, err := bw.Write(input[i : i+1])
		log.PanicIf(err)

		// Seek to somewhere else in the file on the underlying ReadSeeker. It
		// shouldn't affect our reads.
		_, err = sb.Seek(3, os.SEEK_CUR)
		log.PanicIf(err)
	}

	actualBytes := sb.Bytes()
	if bytes.Compare(actualBytes, input) != 0 {
		t.Fatalf("Written buffer not correct: (%d) [%s] %v", len(actualBytes), string(actualBytes), actualBytes)
	}

	seeks := bw.StatsSeeks()
	if seeks != len(input)-1 {
		t.Fatalf("Seek count is not correct (2): (%d)", seeks)
	}

	writes := bw.StatsWrites()
	if writes != len(input) {
		t.Fatalf("Write count is not correct (2): (%d)", writes)
	}
}
