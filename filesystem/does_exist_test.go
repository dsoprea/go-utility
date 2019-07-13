package rifs

import (
	"testing"
)

func TestDoesExist__False(t *testing.T) {
	if DoesExist("a/b/c") != false {
		t.Fatalf("Nonexistent file incorrectly looks like it exists.")
	}
}

func TestDoesExist__True(t *testing.T) {
	if DoesExist("/") != true {
		t.Fatalf("Existent file incorrectly looks like it does not exist.")
	}
}
