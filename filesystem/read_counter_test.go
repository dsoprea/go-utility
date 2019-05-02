package rifs

import (
	"bytes"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-logging"
)

func TestReadCounter_Read(t *testing.T) {
	s := "testing string"
	b := bytes.NewBufferString(s)

	rc := NewReadCounter(b)

	recovered, err := ioutil.ReadAll(rc)
	log.PanicIf(err)

	if bytes.Compare(recovered, []byte(s)) != 0 {
		t.Fatalf("Recovered data not correct:\nACTUAL:\n%v\nEXPECTED:\n%v", recovered, []byte(s))
	}

	if rc.Count() != len(s) {
		t.Fatalf("Counter not correct: (%d)", rc.Count())
	}

	rc.Reset()

	if rc.Count() != 0 {
		t.Fatalf("Count after reset not zero: (%d)", rc.Count())
	}
}
