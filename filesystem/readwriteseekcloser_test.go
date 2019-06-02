package rifs

import (
	"io"
	"testing"
)

func TestReadWriteSeekNoopCloser(t *testing.T) {
	sb := NewSeekableBuffer()

	var c io.Closer
	c = ReadWriteSeekNoopCloser(sb)
	c.Close()
}

func TestReadWriteSeekNoopCloser_Close(t *testing.T) {
	sb := NewSeekableBuffer()

	var c io.Closer
	c = ReadWriteSeekNoopCloser(sb)
	c.Close()
}
