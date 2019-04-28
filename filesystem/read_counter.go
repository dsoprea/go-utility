package rifs

import (
    "io"
)

type ReadCounter struct {
    r       io.Reader
    counter int
}

func NewReadCounter(r io.Reader) *ReadCounter {
    return &ReadCounter{
        r: r,
    }
}

func (rc *ReadCounter) Count() int {
    return rc.counter
}

func (rc *ReadCounter) Reset() {
    rc.counter = 0
}

func (rc *ReadCounter) Read(b []byte) (n int, err error) {
    n, err = rc.r.Read(b)
    rc.counter += n

    return n, err
}
