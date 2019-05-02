package ricrypto

import (
	"hash"
	"io"

	"github.com/dsoprea/go-logging"
)

// ReaderHashProxy proxies a reader and produces a `Hash` sum from the read
// bytes.
type ReaderHashProxy struct {
	r io.Reader
	h hash.Hash
}

// NewReaderHashProxy returns a new `ReaderHashProxy` struct.
func NewReaderHashProxy(r io.Reader, h hash.Hash) *ReaderHashProxy {
	return &ReaderHashProxy{
		r: r,
		h: h,
	}
}

// Read proxies the read to the underlying `Reader` while also pushing the bytes
// through the `Hash` struct.
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

// Sum returns the accumulated hash value.
func (rhp *ReaderHashProxy) Sum() []byte {
	return rhp.h.Sum(nil)
}

// ReaderHash32Proxy proxies a reader and produces a `Hash32` sum from the read
// bytes.
type ReaderHash32Proxy struct {
	r io.Reader
	h hash.Hash32
}

// NewReaderHash32Proxy returns a new `ReaderHash32Proxy` struct.
func NewReaderHash32Proxy(r io.Reader, h hash.Hash32) *ReaderHash32Proxy {
	return &ReaderHash32Proxy{
		r: r,
		h: h,
	}
}

// Read proxies the read to the underlying `Reader` while also pushing the bytes
// through the `Hash32` struct.
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

// Sum32 returns the accumulated hash value.
func (rhp *ReaderHash32Proxy) Sum32() uint32 {
	return rhp.h.Sum32()
}
