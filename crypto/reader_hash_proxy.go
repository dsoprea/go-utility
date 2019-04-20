package ricrypto

import (
    "hash"
    "io"

    "github.com/dsoprea/go-logging"
)

type ReaderHashProxy struct {
    r io.Reader
    h hash.Hash
}

func NewReaderHashProxy(r io.Reader, h hash.Hash) *ReaderHashProxy {
    return &ReaderHashProxy{
        r: r,
        h: h,
    }
}

func (rhp *ReaderHashProxy) Read(b []byte) (n int, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    n, err = rhp.r.Read(b)
    if err != nil {
        if err == io.EOF {
            return 0, err
        }

        log.Panic(err)
    }

    n, err = rhp.h.Write(b[:n])
    log.PanicIf(err)

    return n, nil
}

func (rhp *ReaderHashProxy) Sum() []byte {
    return rhp.h.Sum(nil)
}
