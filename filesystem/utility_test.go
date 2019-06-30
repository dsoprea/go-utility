package rifs

import (
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestGetOffset(t *testing.T) {
	sb := NewSeekableBuffer()

	_, err := sb.Seek(10, os.SEEK_SET)
	log.PanicIf(err)

	n := GetOffset(sb)
	if n != 10 {
		t.Fatalf("Offset not correct.")
	}
}
