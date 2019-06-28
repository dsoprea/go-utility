package rifs

import (
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestBoundedReadWriteSeeker_Seek__SeekCur_Negative(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(-10, os.SEEK_CUR)
	log.PanicIf(err)

	if offsetRaw != 0 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekCur_PositiveThenNegative(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(10, os.SEEK_CUR)
	log.PanicIf(err)

	if offsetRaw != 10 {
		t.Fatalf("Offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+10 {
		t.Fatalf("Offset not correct: (%d)", realOffsetRaw)
	}

	offsetRaw, err = brws.Seek(-20, os.SEEK_CUR)
	log.PanicIf(err)

	if offsetRaw != 0 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw = GetOffset(sb)

	if realOffsetRaw != 55 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekCur_Zero(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(0, os.SEEK_CUR)
	log.PanicIf(err)

	if offsetRaw != 0 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekCur_Positive(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(10, os.SEEK_CUR)
	log.PanicIf(err)

	if offsetRaw != 10 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+10 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekSet_Negative(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(-10, os.SEEK_SET)
	log.PanicIf(err)

	if offsetRaw != 0 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekSet_Zero(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	if offsetRaw != 0 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+0 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekSet_Positive_Small(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(10, os.SEEK_SET)
	log.PanicIf(err)

	if offsetRaw != 10 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+10 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekSet_Positive_Big(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(60, os.SEEK_SET)
	log.PanicIf(err)

	if offsetRaw != 60 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+60 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekEnd_Negative(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(-10, os.SEEK_END)
	log.PanicIf(err)

	if offsetRaw != 99-10 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+99-10 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekEnd_Zero(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(0, os.SEEK_END)
	log.PanicIf(err)

	if offsetRaw != 99 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 55+99 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__SeekEnd_Positive(t *testing.T) {
	sb := NewSeekableBuffer()

	brws, err := NewBoundedReadWriteSeeker(sb, 55, 99)
	log.PanicIf(err)

	_, err = brws.Seek(10, os.SEEK_END)
	if err == nil {
		t.Fatalf("Expected error for seek beyond boundary.")
	} else if err != ErrSeekBeyondBound {
		log.Panic(err)
	}
}

func TestBoundedReadWriteSeeker_Seek__DynamicFileSize__SeekEnd_Negative(t *testing.T) {
	sb := NewSeekableBuffer()

	// Add twenty bytes.
	_, err := sb.Write([]byte("01234567890123456789"))
	log.PanicIf(err)

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	brws, err := NewBoundedReadWriteSeeker(sb, 5, 0)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(-10, os.SEEK_END)
	log.PanicIf(err)

	if offsetRaw != 15-10 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 5+15-10 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__DynamicFileSize__SeekEnd_Zero(t *testing.T) {
	sb := NewSeekableBuffer()

	// Add twenty bytes
	_, err := sb.Write([]byte("01234567890123456789"))
	log.PanicIf(err)

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	brws, err := NewBoundedReadWriteSeeker(sb, 5, 0)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(0, os.SEEK_END)
	log.PanicIf(err)

	if offsetRaw != 15 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 5+15+0 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}

func TestBoundedReadWriteSeeker_Seek__DynamicFileSize__SeekEnd_Positive(t *testing.T) {
	sb := NewSeekableBuffer()

	_, err := sb.Write([]byte("01234567890123456789"))
	log.PanicIf(err)

	_, err = sb.Seek(0, os.SEEK_SET)
	log.PanicIf(err)

	brws, err := NewBoundedReadWriteSeeker(sb, 5, 0)
	log.PanicIf(err)

	offsetRaw, err := brws.Seek(10, os.SEEK_END)
	log.PanicIf(err)

	if offsetRaw != 15+10 {
		t.Fatalf("Relative offset not correct: (%d)", offsetRaw)
	}

	realOffsetRaw := GetOffset(sb)

	if realOffsetRaw != 5+15+10 {
		t.Fatalf("Real offset not correct: (%d)", realOffsetRaw)
	}
}
