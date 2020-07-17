package ritesting

import (
	"testing"

	"github.com/dsoprea/go-logging"
)

func TestGetAndDisableMarshaledExits(t *testing.T) {
	if exitsAreMarshaled != false {
		t.Fatalf("Initial marhsaling state not correct.")
	}

	EnableMarshaledExits()

	if exitsAreMarshaled != true {
		t.Fatalf("Disabled state not correct.")
	}

	defer func() {
		state := recover()
		if state == nil {
			t.Fatalf("main() didn't fail as expected")
		}

		err := state.(error)
		if err.Error() != "marshaled exit of (88)" {
			log.Panic(err)
		}

		DisableMarshaledExits()

		if exitsAreMarshaled != false {
			t.Fatalf("Final marhsaling state not correct.")
		}
	}()

	Exit(88)
}
