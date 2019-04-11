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

    reads := br.StatsReads()
    if reads != 10 {
        t.Fatalf("Read count is not correct: (%d)", reads)
    }

    seeks := br.StatsSeeks()
    if seeks != 9 {
        t.Fatalf("Seek count is not correct: (%d)", seeks)
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

    seeks = br.StatsSeeks()
    if seeks != 11 {
        t.Fatalf("Seek count is not correct (2): (%d)", seeks)
    }

    reads = br.StatsReads()
    if reads != 12 {
        t.Fatalf("Read count is not correct (2): (%d)", reads)
    }
}
