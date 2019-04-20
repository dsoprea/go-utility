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

type ReaderHash32Proxy struct {
    r io.Reader
    h hash.Hash32
}

func NewReaderHash32Proxy(r io.Reader, h hash.Hash32) *ReaderHash32Proxy {
    return &ReaderHash32Proxy{
        r: r,
        h: h,
    }
}

func (rhp *ReaderHash32Proxy) Read(b []byte) (n int, err error) {
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

func (rhp *ReaderHash32Proxy) Sum32() uint32 {
    return rhp.h.Sum32()
}
