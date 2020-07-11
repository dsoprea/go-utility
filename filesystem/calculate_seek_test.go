package rifs

import (
	"io"
	"os"
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestCalculateSeek_SeekCur_Negative_Small(t *testing.T) {
	finalOffset, err := CalculateSeek(11, -8, os.SEEK_CUR, 0)
	log.PanicIf(err)

	if finalOffset != 3 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekCur_Negative_Large(t *testing.T) {
	finalOffset, err := CalculateSeek(11, -20, os.SEEK_CUR, 0)
	log.PanicIf(err)

	if finalOffset != 0 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekCur_Zero(t *testing.T) {
	finalOffset, err := CalculateSeek(11, 0, os.SEEK_CUR, 0)
	log.PanicIf(err)

	if finalOffset != 11 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekCur_Positive(t *testing.T) {
	finalOffset, err := CalculateSeek(11, 22, os.SEEK_CUR, 0)
	log.PanicIf(err)

	if finalOffset != 33 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekSet_Negative(t *testing.T) {
	finalOffset, err := CalculateSeek(11, -8, os.SEEK_SET, 0)
	log.PanicIf(err)

	if finalOffset != 0 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekSet_Zero(t *testing.T) {
	finalOffset, err := CalculateSeek(11, 0, os.SEEK_SET, 0)
	log.PanicIf(err)

	if finalOffset != 0 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekSet_Positive(t *testing.T) {
	finalOffset, err := CalculateSeek(11, 10, os.SEEK_SET, 0)
	log.PanicIf(err)

	if finalOffset != 10 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekEnd_Negative(t *testing.T) {
	finalOffset, err := CalculateSeek(11, -10, os.SEEK_END, 100)
	log.PanicIf(err)

	if finalOffset != 90 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekEnd_Zero(t *testing.T) {
	finalOffset, err := CalculateSeek(11, 0, os.SEEK_END, 100)
	log.PanicIf(err)

	if finalOffset != 100 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek_SeekEnd_Positive(t *testing.T) {
	finalOffset, err := CalculateSeek(11, 10, os.SEEK_END, 100)
	log.PanicIf(err)

	if finalOffset != 110 {
		t.Fatalf("Offset not correct: (%d)", finalOffset)
	}
}

func TestCalculateSeek__BadWhence(t *testing.T) {
	_, err := CalculateSeek(11, 22, 99, 0)
	if err == nil {
		t.Fatalf("Expected failure for bad whence.")
	} else if err.Error() != "whence not valid: (99)" {
		log.Panic(err)
	}
}

func TestSeekType_String_Current(t *testing.T) {
	s := SeekType(io.SeekCurrent).String()
	if s != "SEEK-CURRENT" {
		t.Fatalf("'Current' description not correct: [%s]", s)
	}
}

func TestSeekType_String_Start(t *testing.T) {
	s := SeekType(io.SeekStart).String()
	if s != "SEEK-START" {
		t.Fatalf("'Start' description not correct: [%s]", s)
	}
}

func TestSeekType_String_End(t *testing.T) {
	s := SeekType(io.SeekEnd).String()
	if s != "SEEK-END" {
		t.Fatalf("'End' description not correct: [%s]", s)
	}
}
