package rifs

import (
    "io"
    "os"

    "github.com/dsoprea/go-logging"
)

// BouncebackReader wraps a ReadSeeker, keeps track of our position, and
// seeks back to it before writing. This allows an underlying ReadWriteSeeker
// with an unstable position can still be used for a prolonged series of writes.
type BouncebackReader struct {
    rs              io.ReadSeeker
    currentPosition int64

    statsReads int
    statsSeeks int
}

func NewBouncebackReader(rs io.ReadSeeker) (br *BouncebackReader, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    initialPosition, err := rs.Seek(0, os.SEEK_CUR)
    log.PanicIf(err)

    br = &BouncebackReader{
        rs:              rs,
        currentPosition: initialPosition,
    }

    return br, nil
}

func (br *BouncebackReader) StatsReads() int {
    return br.statsReads
}

func (br *BouncebackReader) StatsSeeks() int {
    return br.statsSeeks
}

func (br *BouncebackReader) Seek(offset int64, whence int) (newPosition int64, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    newPosition, err = br.rs.Seek(offset, whence)
    log.PanicIf(err)

    // Update our internal tracking.
    br.currentPosition = newPosition

    return newPosition, nil
}

func (br *BouncebackReader) Read(p []byte) (n int, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    br.statsReads++

    // Make sure we're where we're supposed to be.

    // This should have no overhead, and enables us to collect stats.
    realCurrentPosition, err := br.rs.Seek(br.currentPosition, os.SEEK_CUR)
    log.PanicIf(err)

    if realCurrentPosition != br.currentPosition {
        br.statsSeeks++

        _, err = br.rs.Seek(br.currentPosition, os.SEEK_SET)
        log.PanicIf(err)
    }

    // Do read.

    n, err = br.rs.Read(p)
    if err != nil {
        if err == io.EOF {
            return 0, io.EOF
        }

        log.Panic(err)
    }

    // Update our internal tracking.
    br.currentPosition += int64(n)

    return n, nil
}
