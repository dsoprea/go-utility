package rifs

import (
    "io"
    "os"

    "github.com/dsoprea/go-logging"
)

// SeekableBuffer is a simple memory structure that satisfies
// `io.ReadWriteSeeker`.
type SeekableBuffer struct {
    data     []byte
    position int64
}

func NewSeekableBuffer() *SeekableBuffer {
    data := make([]byte, 0)

    return &SeekableBuffer{
        data: data,
    }
}

func len64(data []byte) int64 {
    return int64(len(data))
}

func (sb *SeekableBuffer) Len() int {
    return len(sb.data)
}

func (sb *SeekableBuffer) Write(p []byte) (n int, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    // The current position we're already at is past the end of the data we
    // actually have. Extend our buffer up to our current position.
    if sb.position > len64(sb.data) {
        extra := make([]byte, sb.position-len64(sb.data))
        sb.data = append(sb.data, extra...)
    }

    positionFromEnd := len64(sb.data) - sb.position
    tailCount := positionFromEnd - len64(p)

    var tailBytes []byte
    if tailCount > 0 {
        tailBytes = sb.data[len64(sb.data)-tailCount:]
        sb.data = append(sb.data[:sb.position], p...)
    } else {
        sb.data = append(sb.data[:sb.position], p...)
    }

    if tailBytes != nil {
        sb.data = append(sb.data, tailBytes...)
    }

    len_ := len64(p)
    sb.position += len_

    return int(len_), nil
}

func (sb *SeekableBuffer) Read(p []byte) (n int, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if sb.position >= len64(sb.data) {
        return 0, io.EOF
    }

    n = copy(p, sb.data[sb.position:])
    sb.position += int64(n)

    return n, nil
}

func (sb *SeekableBuffer) Seek(offset int64, whence int) (n int64, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    if whence == os.SEEK_SET {
        sb.position = offset
    } else if whence == os.SEEK_END {
        sb.position = len64(sb.data) - offset
    } else if whence == os.SEEK_CUR {
        sb.position += offset
    } else {
        log.Panicf("seek whence is not valid: (%d)", whence)
    }

    return sb.position, nil
}
