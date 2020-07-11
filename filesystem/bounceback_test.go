package rifs

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestBouncebackReader(t *testing.T) {
	originalBytes := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := bytes.NewReader(originalBytes)

	br, err := NewBouncebackReader(r)
	log.PanicIf(err)

	collected := make([]byte, 0)
	for i := 0; i < 10; i++ {
		buffer := make([]byte, 1)

		_, err := br.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Panic(err)
			}
		}

		collected = append(collected, buffer...)

		// Seek to somewhere else in the file on the underlying ReadSeeker. It
		// shouldn't affect our reads.
		_, err = r.Seek(3, os.SEEK_CUR)
		log.PanicIf(err)
	}

	expected := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	if bytes.Compare(collected, expected) != 0 {
		t.Fatalf("Collected bytes not correct.\nACTUAL:\n%v\nEXPECTED:\n%v", collected, expected)
	}

	expectedStats := BouncebackStats{
		reads: 10,
		syncs: 9,
	}

	if br.bouncebackBase.stats != expectedStats {
		t.Fatalf("Stats not correct: %s", br.bouncebackBase.stats)
	}

	// Do another couple of reads that won't actually do anything.

	buffer := make([]byte, 1)

	_, err = br.Read(buffer)
	if err != io.EOF {
		log.PanicIf(err)
	}

	_, err = br.Read(buffer)
	if err != io.EOF {
		log.PanicIf(err)
	}

	// We should've bumped the reads and seeks by two.

	expectedStats = BouncebackStats{
		reads: 12,
		syncs: 10,
	}

	if br.bouncebackBase.stats != expectedStats {
		t.Fatalf("Stats not correct: %s", br.bouncebackBase.stats)
	}
}

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

func TestBouncebackBase_checkPosition__without_sync(t *testing.T) {
	bb := bouncebackBase{
		currentPosition: 0,
	}

	sb := NewSeekableBuffer()

	err := bb.checkPosition(sb)
	log.PanicIf(err)

	actualPosition, err := sb.Seek(0, io.SeekCurrent)
	log.PanicIf(err)

	if actualPosition != 0 {
		t.Fatalf("Actual position is not correct: (%d)", actualPosition)
	}

	expectedStats := BouncebackStats{}

	if bb.stats != expectedStats {
		t.Fatalf("Stats not correct: %s", bb.stats)
	}
}

func TestBouncebackBase_checkPosition__with_sync(t *testing.T) {
	bb := bouncebackBase{
		currentPosition: 11,
	}

	sb := NewSeekableBuffer()

	err := bb.checkPosition(sb)
	log.PanicIf(err)

	actualPosition, err := sb.Seek(0, io.SeekCurrent)
	log.PanicIf(err)

	if actualPosition != 11 {
		t.Fatalf("Actual position is not correct: (%d)", actualPosition)
	}

	expectedStats := BouncebackStats{
		syncs: 1,
	}

	if bb.stats != expectedStats {
		t.Fatalf("Stats not correct: %s", bb.stats)
	}
}
