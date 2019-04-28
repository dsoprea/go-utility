package rifs

import (
    "io"
)

type WriteCounter struct {
    w       io.Writer
    counter int
}

func NewWriteCounter(w io.Writer) *WriteCounter {
    return &WriteCounter{
        w: w,
    }
}

func (wc *WriteCounter) Count() int {
    return wc.counter
}

func (wc *WriteCounter) Reset() {
    wc.counter = 0
}

func (wc *WriteCounter) Write(b []byte) (n int, err error) {
    n, err = wc.w.Write(b)
    wc.counter += n

    return n, err
}
