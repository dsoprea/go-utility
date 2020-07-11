package rifs

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dsoprea/go-logging"
)

func TestNewWriteProgressWrapper(t *testing.T) {
	b := new(bytes.Buffer)
	w := NewWriteProgressWrapper(b, nil)
	wpw := w.(*WriteProgressWrapper)

	if wpw.w != b {
		t.Fatalf("Internal writer not correct.")
	}
}

func TestWriteProgressWrapper_Write(t *testing.T) {
	steps := make([]int, 0)
	progressCb := func(n int, duration time.Duration, isEof bool) error {
		steps = append(steps, n)
		return nil
	}

	b := new(bytes.Buffer)
	wpw := NewWriteProgressWrapper(b, progressCb)

	_, err := wpw.Write([]byte("abc"))
	log.PanicIf(err)

	_, err = wpw.Write([]byte("defg"))
	log.PanicIf(err)

	_, err = wpw.Write([]byte("hijkl"))
	log.PanicIf(err)

	expected := []int{
		3,
		4,
		5,
	}

	if reflect.DeepEqual(steps, expected) != true {
		t.Fatalf("Iterative progress not correct.")
	}
}

func TestNewReadProgressWrapper(t *testing.T) {
	b := new(bytes.Buffer)
	r := NewReadProgressWrapper(b, nil)

	rpw := r.(*ReadProgressWrapper)

	if rpw.r != b {
		t.Fatalf("Internal reader not correct.")
	}
}

func TestReadProgressWrapper_Read(t *testing.T) {
	steps := make([][]interface{}, 0)
	progressCb := func(n int, duration time.Duration, isEof bool) error {
		steps = append(steps, []interface{}{n, isEof})
		return nil
	}

	data := []byte(strings.Repeat("1234567890", 500))
	inputReader := bytes.NewBuffer(data)

	rpw := NewReadProgressWrapper(inputReader, progressCb)

	outputWriter := new(bytes.Buffer)

	_, err := io.Copy(outputWriter, rpw)
	log.PanicIf(err)

	outputBytes := outputWriter.Bytes()

	expectedSteps := [][]interface{}{
		{512, false},
		{1024, false},
		{2048, false},
		{1416, false},
		{0, true},
	}

	if bytes.Equal(outputBytes, data) != true {
		t.Fatalf("Recovered bytes not correct.")
	}

	if reflect.DeepEqual(steps, expectedSteps) != true {
		t.Fatalf("Steps not correct.")
	}
}
